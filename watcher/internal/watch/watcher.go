package watch

import (
	"context"
	"strings"
	"time"

	"github.com/anton.okolelov/pgquerywatcher/internal/config"
	"github.com/anton.okolelov/pgquerywatcher/internal/pg"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Watcher struct {
	targetDb *pgxpool.Pool
	log      zerolog.Logger
}

type PgStatStatementsRow struct {
	Query     string
	Calls     uint64
	TotalTime float64
}

type QueryStats struct {
	Time           time.Time
	Query          string
	ExecMeanTimeMs float64
}

func NewWatcher(targetDbConfig config.DbConfig, log zerolog.Logger) (*Watcher, error) {
	targetDb, err := pg.NewDB(targetDbConfig)
	if err != nil {
		return nil, err
	}

	return &Watcher{targetDb: targetDb, log: log}, err
}

func (w Watcher) Watch() error {
	previousStatStatements, err := w.getStatStatements()
	if err != nil {
		return err
	}

	for range time.NewTicker(3 * time.Second).C {
		currentStatStatements, err := w.getStatStatements()
		if err != nil {
			return err
		}

		w.logChanges(currentStatStatements, previousStatStatements)
		previousStatStatements = currentStatStatements
	}
	return nil
}

func (w Watcher) logChanges(
	currentStatStatements map[string]PgStatStatementsRow,
	previousStatStatements map[string]PgStatStatementsRow,
) {
	currTime := time.Now()

	for query, curr := range currentStatStatements {
		if strings.Contains(query, "pg_stat_statements") {
			continue
		}
		prev := previousStatStatements[query]

		// if someone reset stats, ignore previous values
		if prev.Calls > curr.Calls {
			prev = PgStatStatementsRow{}
		}

		if curr.Calls == prev.Calls {
			continue
		}

		timeDiff := curr.TotalTime - prev.TotalTime
		countDiff := curr.Calls - prev.Calls
		execMeanTime := timeDiff / float64(countDiff)

		w.log.Info().
			Bool("is_sql_stats", true).
			Time("time", currTime).
			Float64("time_diff", timeDiff).
			Uint64("count_diff", countDiff).
			Float64("exec_mean_time", execMeanTime).
			Msg(query)
	}
}

func (w Watcher) getStatStatements() (map[string]PgStatStatementsRow, error) {

	rows, err := w.targetDb.Query(context.Background(), `
		SELECT query, calls, total_exec_time 
			FROM pg_stat_statements`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	queriesStats := make(map[string]PgStatStatementsRow)

	for rows.Next() {
		var queryStat PgStatStatementsRow
		err = rows.Scan(&queryStat.Query, &queryStat.Calls, &queryStat.TotalTime)
		if err != nil {
			return nil, err
		}
		queriesStats[queryStat.Query] = queryStat
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return queriesStats, nil
}
