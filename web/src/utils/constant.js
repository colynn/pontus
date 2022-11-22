export const ENVS = [
  {
    value: 'prod',
    label: '生产环境'
  },
  {
    value: 'uat',
    label: 'UAT环境'
  },
  {
    value: 'event',
    label: '事件环境'
  },
  {
    value: 'demo',
    label: '演示环境'
  },
  {
    value: 'public',
    label: '共用环境'
  }

]

export const Highlight = 'alter|and|as|asc|between|count|create|delete|desc|distinct|drop|from|having|in|insert|into|is|join|like|not|on|or|order|select|set|table|union|update|values|where|accessible|action|add|after|algorithm|all|analyze|asensitive|at|authors|auto_increment|autocommit|avg|avg_row_length|before|binary|binlog|both|btree|cache|call|cascade|cascaded|case|catalog_name|chain|change|changed|character|check|checkpoint|checksum|class_origin|client_statistics|close|coalesce|code|collate|collation|collations|column|columns|comment|commit|committed|completion|concurrent|condition|connection|consistent|constraint|contains|continue|contributors|convert|cross|current_date|current_time|current_timestamp|current_user|cursor|data|database|databases|day_hour|day_microsecond|day_minute|day_second|deallocate|dec|declare|default|delay_key_write|delayed|delimiter|des_key_file|describe|deterministic|dev_pop|dev_samp|deviance|directory|disable|discard|distinctrow|div|dual|dumpfile|each|elseif|enable|enclosed|end|ends|engine|engines|enum|errors|escape|escaped|even|event|events|every|execute|exists|exit|explain|extended|fast|fetch|field|fields|first|flush|for|force|foreign|found_rows|full|fulltext|function|general|global|grant|grants|group|groupby_concat|handler|hash|help|high_priority|hosts|hour_microsecond|hour_minute|hour_second|if|ignore|ignore_server_ids|import|index|index_statistics|infile|inner|innodb|inout|insensitive|insert_method|install|interval|invoker|isolation|iterate|key|keys|kill|language|last|leading|leave|left|level|limit|linear|lines|list|load|local|localtime|localtimestamp|lock|logs|low_priority|master|master_heartbeat_period|master_ssl_verify_server_cert|masters|match|max|max_rows|maxvalue|message_text|middleint|migrate|min|min_rows|minute_microsecond|minute_second|mod|mode|modifies|modify|mutex|mysql_errno|natural|next|no|no_write_to_binlog|offline|offset|one|online|open|optimize|option|optionally|out|outer|outfile|pack_keys|parser|partition|partitions|password|phase|plugin|plugins|prepare|preserve|prev|primary|privileges|procedure|processlist|profile|profiles|purge|query|quick|range|read|read_write|reads|real|rebuild|recover|references|regexp|relaylog|release|remove|rename|reorganize|repair|repeatable|replace|acequire|resignal|restrict|resume|return|returns|revoke|right|rlike|rollback|rollup|row|row_format|rtree|savepoint|schedule|schema|schema_name|schemas|second_microsecond|security|sensitive|separator|serializable|server|session|share|show|signal|slave|slow|smallint|snapshot|soname|spatial|specific|sql|sql_big_result|sql_buffer_result|sql_cache|sql_calc_found_rows|sql_no_cache|sql_small_result|sqlexception|sqlstate|sqlwarning|ssl|start|starting|starts|status|std|stddev|stddev_pop|stddev_samp|storage|straight_join|subclass_origin|sum|suspend|table_name|table_statistics|tables|tablespace|temporary|terminated|to|trailing|transaction|trigger|triggers|truncate|uncommitted|undo|uninstall|unique|unlock|upgrade|usage|use|use_frm|user|user_resources|user_statistics|using|utc_date|utc_time|utc_timestamp|value|variables|varying|view|views|warnings|when|while|with|work|write|xa|xor|year_month|zerofill|begin|do|then|else|loop|repeat|by|bool|boolean|bit|blob|decimal|double|enum|float|long|longblob|longtext|medium|mediumblob|mediumint|mediumtext|time|timestamp|tinyblob|tinyint|tinytext|text|bigint|int|int1|int2|int3|int4|int8|integer|float|float4|float8|double|char|varbinary|varchar|varcharacter|precision|date|datetime|year|unsigned|signed|numeric|ucase|lcase|mid|len|round|rank|now|format|coalesce|ifnull|isnull|nvl'

export const AssetStatus = [
  {
    value: 'Running',
    label: '在线'
  },
  {
    value: 'Stopped',
    label: '关闭'
  }
  // {
  //   value: 'OffLine',
  //   label: '下线'
  // }
]

export const OSType = [
  {
    value: 'linux',
    label: 'linux'
  },
  {
    value: 'windows',
    label: 'windows'
  }
]

export const AssetType = [
  {
    value: 'ecs',
    label: '云服务器 ECS'
  },
  {
    value: 'slb',
    label: '负载均衡'
  },
  {
    value: 'nas',
    label: '文件存储NAS'
  },
  {
    value: 'oss',
    label: '对象存储'
  },
  {
    value: 'eip',
    label: '弹性公网IP'
  },
  {
    value: 'nats',
    label: 'NAT网关'
  },
  {
    value: 'disk',
    label: '块存储'
  },
  {
    value: 'ga',
    label: '全球加速'
  },
  {
    value: 'cdn',
    label: 'CDN'
  },
  {
    value: 'waf',
    label: 'Web应用防火墙'
  }
]

export const TangibleAssetStatus = [
  {
    value: 1,
    label: '使用中'
  },
  {
    value: 2,
    label: '闲置中'
  },
  {
    value: 3,
    label: '需报废'
  },
  {
    value: 4,
    label: '已报废'
  },
  {
    value: 5,
    label: '需退租'
  },
  {
    value: 6,
    label: '已退租'
  }
]
