package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	pkgErrors "github.com/pkg/errors"
	"github.com/prestodb/presto-go-client/presto"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

const DefaultQuery = "SELECT 1, 2, 3"

type PrestoDatasource struct {
	settings *DatasourceSettings
	db       *sql.DB
}

type DatasourceSettings struct {
	Instance    backend.DataSourceInstanceSettings
	PrestoParam PrestoParam
}

type PrestoParam struct {
	HTTPScheme   string
	Host         string
	Catalog      string
	Schema       string
	CustomParams []struct {
		Name  string
		Value string
	}
	RowLimit                 int64
	ResultRowLimit           int64
	QueryMaxExecutionSeconds int64
}

type Query struct {
	RefId  string `json:"refId"`
	RawSql string `json:"rawSql"`
	Format string `json:"format"`
}

func NewDatasourceInstance(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var dsSettings = DatasourceSettings{
		Instance: settings,
	}
	if err := json.Unmarshal(dsSettings.Instance.JSONData, &dsSettings.PrestoParam); err != nil {
		return nil, fmt.Errorf("unable to parse settings json %s. Error: %w", dsSettings.Instance.JSONData, err)
	}
	if dsSettings.PrestoParam.RowLimit <= 0 {
		dsSettings.PrestoParam.RowLimit = 1000000
	}
	if dsSettings.PrestoParam.ResultRowLimit < 0 {
		dsSettings.PrestoParam.ResultRowLimit = 0
	}
	if dsSettings.PrestoParam.QueryMaxExecutionSeconds <= 0 {
		dsSettings.PrestoParam.QueryMaxExecutionSeconds = 60
	}
	dsn := fmt.Sprintf("%s://%s@%s?catalog=%s&schema=%s&custom_client=%s&session_properties=query_max_execution_time=%ds",
		dsSettings.PrestoParam.HTTPScheme,
		dsSettings.Instance.BasicAuthUser,
		dsSettings.PrestoParam.Host,
		dsSettings.PrestoParam.Catalog,
		dsSettings.PrestoParam.Schema,
		dsSettings.Instance.Name,
		dsSettings.PrestoParam.QueryMaxExecutionSeconds,
	)
	params := url.Values{}
	for _, param := range dsSettings.PrestoParam.CustomParams {
		params.Add(param.Name, param.Value)
	}
	if len(params) > 0 {
		dsn = dsn + "&" + params.Encode()
	}
	db, err := sql.Open("presto", dsn)
	if err != nil {
		return nil, err
	}
	presto.RegisterCustomClient(dsSettings.Instance.Name, &http.Client{})
	backend.Logger.Info("Create datasource.", "datasource", dsSettings.Instance.Name, "url", dsn)
	return &PrestoDatasource{
		settings: &dsSettings,
		db:       db,
	}, nil
}

func (ds *PrestoDatasource) Dispose() {
	backend.Logger.Info("Dispose datasource.", "datasource", ds.settings.Instance.Name)
	ds.db.Close()
	presto.DeregisterCustomClient(ds.settings.Instance.Name)
}

func (ds *PrestoDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	onErr := func(err error) (*backend.QueryDataResponse, error) {
		backend.Logger.Error(fmt.Sprintf("QueryData error: %v", err))
		return nil, err
	}
	result := backend.NewQueryDataResponse()
	ch := make(chan DBDataResponse, len(req.Queries))
	var wg sync.WaitGroup
	for _, query := range req.Queries {
		var queryJson = Query{}
		err := json.Unmarshal(query.JSON, &queryJson)
		if err != nil {
			return onErr(fmt.Errorf("unable to parse json %s. Error: %w", query.JSON, err))
		}
		wg.Add(1)
		go ds.queryData(query, &wg, ctx, ch, queryJson)
	}
	wg.Wait()
	close(ch)
	result.Responses = make(map[string]backend.DataResponse)
	for queryResult := range ch {
		result.Responses[queryResult.refID] = queryResult.dataResponse
	}
	return result, nil
}

func (ds *PrestoDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	onErr := func(err error) (*backend.CheckHealthResult, error) {
		backend.Logger.Error(fmt.Sprintf("HealthCheck error: %v", err))
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, err
	}
	backend.Logger.Info(fmt.Sprintf("Starting HealthCheck, req:%v", req))

	rows, err := ds.queryRows(DefaultQuery)
	if err != nil {
		return onErr(err)
	}
	rows.Close()

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "OK",
	}, nil
}

func (ds *PrestoDatasource) queryData(query backend.DataQuery, wg *sync.WaitGroup, queryContext context.Context,
	ch chan DBDataResponse, queryJson Query) {
	defer wg.Done()
	queryResult := DBDataResponse{
		dataResponse: backend.DataResponse{},
		refID:        query.RefID,
	}
	onErr := func(err error) {
		queryResult.dataResponse.Error = err
		queryError.Add(1)
		backend.Logger.Error(fmt.Sprintf("Datasource query error: %s", err))
	}
	defer func(start time.Time) {
		if r := recover(); r != nil {
			backend.Logger.Error("executeQuery panic", "error", r)
			queryResult.dataResponse.Error = fmt.Errorf("%v", r)
		}
		queryPrestoCost.Observe(float64(time.Since(start).Milliseconds()))
		ch <- queryResult
	}(time.Now())
	rows, err := ds.queryRows(queryJson.RawSql)
	if err != nil {
		onErr(err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			backend.Logger.Warn("Failed to close rows.", "err", err.Error())
		}
	}()

	qm, err := newProcessCfg(query, queryContext, rows)
	if err != nil {
		onErr(err)
		return
	}

	// Convert row.Rows to dataframe
	frame, err := sqlutil.FrameFromRows(rows, ds.settings.PrestoParam.RowLimit, Converters()...)
	if err != nil {
		onErr(err)
		return
	}

	if frame.Meta == nil {
		frame.Meta = &data.FrameMeta{}
	}

	frame.Meta.ExecutedQueryString = queryJson.RawSql

	// If no rows were returned, no point checking anything else.
	if frame.Rows() == 0 {
		queryResult.dataResponse.Frames = data.Frames{frame}
		return
	}

	if err := convertSQLTimeColumnsToEpochMS(frame, qm); err != nil {
		onErr(err)
		return
	}

	if qm.Format == dataQueryFormatSeries {
		// time series has to have time column
		if qm.timeIndex == -1 {
			onErr(errors.New("no time column found"))
			return
		}

		// Make sure to name the time field 'Time' to be backward compatible with Grafana pre-v8.
		frame.Fields[qm.timeIndex].Name = data.TimeSeriesTimeFieldName

		for i := range qm.columnNames {
			if i == qm.timeIndex || i == qm.metricIndex {
				continue
			}

			if t := frame.Fields[i].Type(); t == data.FieldTypeString || t == data.FieldTypeNullableString {
				continue
			}

			var err error
			if frame, err = convertSQLValueColumnToFloat(frame, i); err != nil {
				onErr(pkgErrors.Wrap(err, "convert value to float failed"))
				return
			}
		}

		tsSchema := frame.TimeSeriesSchema()
		if tsSchema.Type == data.TimeSeriesTypeLong {
			var err error
			originalData := frame
			frame, err = data.LongToWide(frame, qm.FillMissing)
			if err != nil {
				pkgErrors.Wrap(err, "failed to convert long to wide series when converting from dataframe")
				return
			}

			// Before 8x, a special metric column was used to name time series. The LongToWide transforms that into a metric label on the value field.
			// But that makes series name have both the value column name AND the metric name. So here we are removing the metric label here and moving it to the
			// field name to get the same naming for the series as pre v8
			if len(originalData.Fields) == 3 {
				for _, field := range frame.Fields {
					if len(field.Labels) == 1 { // 7x only supported one label
						name, ok := field.Labels["metric"]
						if ok {
							field.Name = name
							field.Labels = nil
						}
					}
				}
			}
		}
		if qm.FillMissing != nil {
			var err error
			frame, err = resample(frame, *qm)
			if err != nil {
				backend.Logger.Error("Failed to resample dataframe", "err", err)
				frame.AppendNotices(data.Notice{Text: "Failed to resample dataframe", Severity: data.NoticeSeverityWarning})
			}
		}
	}
	queryResult.dataResponse.Frames = data.Frames{frame}
}

func (ds *PrestoDatasource) queryRows(query string) (*sql.Rows, error) {
	onErr := func(err error) (*sql.Rows, error) {
		backend.Logger.Error(fmt.Sprintf("presto client query error: %v", err))
		return nil, err
	}
	if strings.TrimSpace(query) != "" && ds.settings.PrestoParam.ResultRowLimit > 0 {
		query = fmt.Sprintf("SELECT * FROM ( %s ) query_limit_wrapper limit %d", query, ds.settings.PrestoParam.ResultRowLimit)
	}
	backend.Logger.Info("Query presto.", "datasource", ds.settings.Instance.Name, "query", query)
	rows, err := ds.db.Query(query)
	if err != nil {
		return onErr(fmt.Errorf("do presto query failed, err: %s", err.Error()))
	}
	return rows, nil
}
