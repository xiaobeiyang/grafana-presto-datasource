import {
  DataSourceInstanceSettings,
  MetricFindValue,
  DataFrame,
  ScopedVars,
  DataQueryResponse,
  DataQueryRequest,
  Field,
} from '@grafana/data';
import {
  BackendDataSourceResponse,
  FetchResponse,
  DataSourceWithBackend,
  toDataQueryResponse,
  getTemplateSrv,
  getBackendSrv,
} from '@grafana/runtime';
import { PrestoDataSourceOptions, PrestoQuery } from './types';
import { map, catchError } from 'rxjs/operators';
import { lastValueFrom, of, Observable } from 'rxjs';
import { each } from 'lodash';

export const FORMAT_TIME_SERIES = 'time_series';
export const FORMAT_TABLE = 'table';

export class DataSource extends DataSourceWithBackend<PrestoQuery, PrestoDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<PrestoDataSourceOptions>) {
    super(instanceSettings);
  }

  filterQuery(query: PrestoQuery): boolean {
    return !query.hide;
  }

  query(request: DataQueryRequest<PrestoQuery>): Observable<DataQueryResponse> {
    each(request.targets, (query) => {
      migrateQuery(query);
    });
    const observableResponse = super.query(request);
    return observableResponse.pipe(
      map((response) => {
        each(response.data, (frame: DataFrame) => {
          let legendFormat = '';
          each(request.targets, (target) => {
            if (target.refId === frame.refId && target.format === FORMAT_TIME_SERIES && target.legendFormat) {
              legendFormat = target.legendFormat;
              return false;
            }
            return;
          });
          if (legendFormat === '') {
            return;
          }
          each(frame.fields, (field: Field) => {
            field.config.displayNameFromDS = this.renderTemplate(
              getTemplateSrv().replace(legendFormat, request.scopedVars),
              field.labels || {}
            );
          });
        });
        return {
          ...response,
        };
      })
    );
  }

  renderTemplate(aliasPattern: string, aliasData: { [key: string]: string }) {
    const aliasRegex = /\{\{\s*(.+?)\s*\}\}/g;
    return aliasPattern.replace(aliasRegex, (_match, g1) => {
      if (aliasData[g1]) {
        return aliasData[g1];
      }
      return '';
    });
  }

  applyTemplateVariables(query: PrestoQuery, scopedVars: ScopedVars) {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      rawSql: query.rawSql ? templateSrv.replace(query.rawSql, scopedVars) : '',
    };
  }

  metricFindQuery(query: string, optionalOptions?: any): Promise<MetricFindValue[]> {
    let refId = 'tempvar';
    if (optionalOptions && optionalOptions.variable && optionalOptions.variable.name) {
      refId = optionalOptions.variable.name;
    }
    let sql = getTemplateSrv().replace(query, optionalOptions.scopedVars);
    return lastValueFrom(
      getBackendSrv()
        .fetch<BackendDataSourceResponse>({
          url: '/api/ds/query',
          method: 'POST',
          data: {
            queries: [
              {
                datasourceId: this.id,
                refId: 'metricFindQuery',
                rawSql: sql,
                format: 'table',
              },
            ],
          },
          requestId: refId,
        })
        .pipe(
          map((rsp) => {
            return this.transformMetricFindResponse(rsp);
          }),
          catchError((err) => {
            return of([]);
          })
        )
    );
  }

  transformMetricFindResponse(raw: FetchResponse<BackendDataSourceResponse>): MetricFindValue[] {
    const frames = toDataQueryResponse(raw).data as DataFrame[];

    if (!frames || !frames.length) {
      return [];
    }

    const frame = frames[0];

    const values: MetricFindValue[] = [];
    const textField = frame.fields.find((f) => f.name === '__text');
    const valueField = frame.fields.find((f) => f.name === '__value');

    if (textField && valueField) {
      for (let i = 0; i < textField.values.length; i++) {
        values.push({ text: '' + textField.values.get(i), value: '' + valueField.values.get(i) });
      }
    } else {
      values.push(
        ...frame.fields
          .flatMap((f) => f.values.toArray())
          .map((v) => ({
            text: v,
          }))
      );
    }

    return Array.from(new Set(values.map((v) => v.text))).map((text) => ({
      text,
      value: values.find((v) => v.text === text)?.value,
    }));
  }
}

// for backward compatibility
export function migrateQuery(query: PrestoQuery) {
  if (query.queryText) {
    query.rawSql = query.queryText;
    if (query.queryType === 'timeseries') {
      query.format = FORMAT_TIME_SERIES;
    } else {
      query.format = FORMAT_TABLE;
    }
    delete query.queryText;
    delete query.queryType;
  }
}
