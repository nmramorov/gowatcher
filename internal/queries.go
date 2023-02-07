package metrics

const (
	CREATE_GAUGE_TABLE string = `CREATE TABLE IF NOT EXISTS gaugeMetrics (
		_id TEXT,
		mtype TEXT,
		_value DOUBLE PRECISION,
	  	date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	CREATE_COUNTER_TABLE string = `CREATE TABLE IF NOT EXISTS counterMetrics (
		_id TEXT,
		mtype TEXT,
		_value INTEGER,
		date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	INSERT_INTO_GAUGE              = `INSERT INTO gaugemetrics VALUES ($1, $2, $3);`
	INSERT_INTO_COUNTER            = `INSERT INTO countermetrics VALUES ($1, $2, $3);`
	SELECT_FROM_GAUGE          = `SELECT * FROM gaugemetrics WHERE _id=$1 ORDER BY date DESC LIMIT 1`
	SELECT_FROM_COUNTER string = `SELECT * FROM countermetrics WHERE _id=$1 ORDER BY date DESC LIMIT 1`
)
