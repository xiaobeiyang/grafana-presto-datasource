import defaults from 'lodash/defaults';

import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select, InlineFormLabel } from '@grafana/ui';
import { DataSource, migrateQuery, FORMAT_TABLE, FORMAT_TIME_SERIES } from './DataSource';
import { defaultQuery, PrestoDataSourceOptions, PrestoQuery } from './types';

import AceEditor from 'react-ace';
import 'ace-builds/src-min-noconflict/ext-language_tools';
import 'ace-builds/src-noconflict/mode-mysql';
import 'ace-builds/src-noconflict/theme-terminal';

type Props = QueryEditorProps<DataSource, PrestoQuery, PrestoDataSourceOptions>;

const FORMAT_OPTIONS: Array<SelectableValue<string>> = [
  { label: 'Time series', value: FORMAT_TIME_SERIES },
  { label: 'Table', value: FORMAT_TABLE },
];

export class QueryEditor extends PureComponent<Props> {
  onQueryChange = (rawSql: string) => {
    const { onChange, query } = this.props;
    onChange({ ...query, rawSql: rawSql, format: query.format || FORMAT_TIME_SERIES });
  };

  onQueryBlur = () => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, format: query.format || FORMAT_TIME_SERIES });
    onRunQuery();
  };

  onFormatChange = (option: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, format: option.value as any });
  };

  onLegendFormatChange = (e: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, legendFormat: e.currentTarget.value });
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    migrateQuery(query);
    const { rawSql, format, legendFormat } = query;
    return (
      <div>
        <div className="gf-form">
          <AceEditor
            placeholder=""
            mode="mysql"
            name="qEditor"
            onChange={this.onQueryChange}
            onBlur={this.onQueryBlur}
            fontSize={14}
            height="200px"
            width="100%"
            showPrintMargin={false}
            showGutter={true}
            highlightActiveLine={true}
            value={rawSql || ''}
            setOptions={{
              enableBasicAutocompletion: true,
              enableLiveAutocompletion: true,
              enableSnippets: false,
              showLineNumbers: true,
              tabSize: 2,
            }}
          />
        </div>
        <div className="gf-form">
          <InlineFormLabel
            className="gf-form-label width-7"
            tooltip="Controls the name of the time series, using name or pattern. For example
        {{hostname}} will be replaced with field value for the field hostname."
          >
            Legend Format
          </InlineFormLabel>
          <input
            type="text"
            className="gf-form-input width-24"
            placeholder="legend format"
            value={legendFormat}
            onChange={this.onLegendFormatChange}
            onBlur={this.onQueryBlur}
          />
          <div className="gf-form-label width-7">Format</div>
          <Select
            menuShouldPortal
            className="select-container"
            width={16}
            isSearchable={false}
            options={FORMAT_OPTIONS}
            onChange={this.onFormatChange}
            onBlur={this.onQueryBlur}
            value={format}
          />
        </div>
      </div>
    );
  }
}
