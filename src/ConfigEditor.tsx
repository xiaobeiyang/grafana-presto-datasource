import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms, Button, Icon } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { PrestoDataSourceOptions } from './types';
import { map, filter } from 'lodash';
const { FormField } = LegacyForms;
interface Props extends DataSourcePluginOptionsEditorProps<PrestoDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onHTTPSchemeChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      httpScheme: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };
  onHostChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      host: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };
  onCatalogChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      catalog: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };
  // Secure field (only sent to the backend)
  onSchemaChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      schema: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };
  onUserChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({ ...options, basicAuthUser: event.currentTarget.value });
  };
  onQueryMaxExecutionSecondsChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    let queryMaxExecutionSeconds = Number(event.target.value);
    if (isNaN(queryMaxExecutionSeconds) || queryMaxExecutionSeconds <= 0) {
      queryMaxExecutionSeconds = 60;
    }
    const jsonData = {
      ...options.jsonData,
      queryMaxExecutionSeconds: queryMaxExecutionSeconds,
    };
    onOptionsChange({ ...options, jsonData });
  };
  onRowLimitChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    let rowLimit = Number(event.target.value);
    if (isNaN(rowLimit) || rowLimit <= 0) {
      rowLimit = 1000000;
    }
    const jsonData = {
      ...options.jsonData,
      rowLimit: rowLimit,
    };
    onOptionsChange({ ...options, jsonData });
  };
  onResultRowLimitChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    let resultRowLimit = Number(event.target.value);
    if (isNaN(resultRowLimit) || resultRowLimit < 0) {
      resultRowLimit = 100000;
    }
    const jsonData = {
      ...options.jsonData,
      resultRowLimit: resultRowLimit,
    };
    onOptionsChange({ ...options, jsonData });
  };

  render() {
    const { options } = this.props;
    const { jsonData } = options;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="HTTP Scheme"
            labelWidth={10}
            inputWidth={30}
            onChange={this.onHTTPSchemeChange}
            value={jsonData.httpScheme || ''}
            placeholder="https or http"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Presto host"
            labelWidth={10}
            inputWidth={30}
            onChange={this.onHostChange}
            value={jsonData.host || ''}
            placeholder="Presto cluster host"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="User"
            labelWidth={10}
            inputWidth={30}
            onChange={this.onUserChange}
            value={options.basicAuthUser}
            placeholder="user"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Catalog"
            labelWidth={10}
            inputWidth={30}
            onChange={this.onCatalogChange}
            value={jsonData.catalog || ''}
            placeholder="A Presto Database Catalog"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Schema"
            labelWidth={10}
            inputWidth={30}
            onChange={this.onSchemaChange}
            value={jsonData.schema || ''}
            placeholder="A Presto Database Schema"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Query max execution seconds"
            type="number"
            labelWidth={20}
            inputWidth={20}
            onChange={this.onQueryMaxExecutionSecondsChange}
            value={jsonData.queryMaxExecutionSeconds}
            placeholder="60"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Row limit"
            type="number"
            labelWidth={10}
            inputWidth={30}
            onChange={this.onRowLimitChange}
            value={jsonData.rowLimit}
            placeholder="1000000"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Result row limit"
            type="number"
            labelWidth={10}
            inputWidth={30}
            onChange={this.onResultRowLimitChange}
            value={jsonData.resultRowLimit}
            placeholder="max result rows, default is 0(no limit)."
          />
        </div>
        <div>
          <CustomUrlParamSettings {...this.props} />
        </div>
      </div>
    );
  }
}

export class CustomUrlParamSettings extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props);
    const { jsonData } = this.props.options;
    this.state = {
      headers: map(jsonData.customParams, (header) => {
        return { ...header };
      }),
    };
  }

  onAdd = () => {
    const { options } = this.props;
    const customParams = this.props.options.jsonData.customParams ? this.props.options.jsonData.customParams : [];
    const newCustomParams = [...customParams, { name: '', value: '' }];
    this.props.onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        customParams: newCustomParams,
      },
    });
  };

  onRemove = (idx: Number) => {
    const { options } = this.props;
    const { customParams } = this.props.options.jsonData;
    const newCustomParams = filter(customParams, (header, i: Number) => i !== idx);
    this.props.onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        customParams: newCustomParams,
      },
    });
  };

  onNameChange = (idx: Number, value: string) => {
    const { options } = this.props;
    const { customParams } = this.props.options.jsonData;
    const newCustomParams = map(customParams, (header, i: Number) => {
      if (i !== idx) {
        return header;
      }
      return {
        ...header,
        name: value,
      };
    });
    this.props.onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        customParams: newCustomParams,
      },
    });
  };

  onValueChange = (idx: Number, value: string) => {
    const { options } = this.props;
    const { customParams } = this.props.options.jsonData;
    const newCustomParams = map(customParams, (header, i: Number) => {
      if (i !== idx) {
        return header;
      }
      return {
        ...header,
        value: value,
      };
    });
    this.props.onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        customParams: newCustomParams,
      },
    });
  };

  render() {
    const { customParams } = this.props.options.jsonData;
    return (
      <div className={'gf-form-group'}>
        <div className="gf-form">
          <h6>Custom HTTP URL Params</h6>
        </div>
        <div>
          {map(customParams, (param, idx: Number) => (
            <div className={'gf-form'}>
              <FormField
                label="Name"
                name="name"
                placeholder="name"
                labelWidth={5}
                value={param.name || ''}
                onChange={(e) => this.onNameChange(idx, e.target.value)}
              />
              <FormField
                label="Value"
                name="value"
                placeholder="value"
                labelWidth={5}
                value={param.value || ''}
                onChange={(e) => this.onValueChange(idx, e.target.value)}
              />
              <Button
                type="button"
                aria-label="Remove"
                variant="secondary"
                size="xs"
                onClick={(_e) => this.onRemove(idx)}
              >
                <Icon name="trash-alt" />
              </Button>
            </div>
          ))}
        </div>
        <div className="gf-form">
          <Button
            variant="secondary"
            icon="plus"
            type="button"
            onClick={(e) => {
              this.onAdd();
            }}
          >
            Add
          </Button>
        </div>
      </div>
    );
  }
}
