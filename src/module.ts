import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { PrestoQuery, PrestoDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, PrestoQuery, PrestoDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
