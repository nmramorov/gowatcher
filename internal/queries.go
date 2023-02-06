package metrics

const (
	CREATE_GAUGE_TABLE string = `CREATE TABLE gaugeMetrics (
		_id TEXT,
		mtype TEXT,
		_value DOUBLE PRECISION
	);`
	CREATE_COUNTER_TABLE string = `CREATE TABLE counterMetrics (
		_id TEXT,
		mtype TEXT,
		_value INTEGER
	);`
	INSERT_INTO_GAUGE              = `INSERT INTO gaugemetrics VALUES ($1, $2, $3);`
	INSERT_INTO_COUNTER            = `INSERT INTO countermetrics VALUES ($1, $2, $3);`
	SELECT_ALL_FROM_GAUGE          = `SELECT * FROM gaugemetrics;`
	SELECT_ALL_FROM_COUNTER string = `SELECT * FROM countermetrics;`
)
