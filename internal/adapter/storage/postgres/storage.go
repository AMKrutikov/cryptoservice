package postgres

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AMKrutikov/cryptoservice/internal/cases"
	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var (
	_ cases.CryptoStorage = (*Storage)(nil)
)

type Storage struct {
	pool     *pgxpool.Pool
	cancelFn context.CancelFunc
	once     sync.Once
}

func NewStorage(connString string) (*Storage, error) {
	ctx, cancel := context.WithCancel(context.Background())

	pgxPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		cancel()
		return nil, errors.Wrapf(entities.ErrInternal, "failed to connect to the database: %v", err)
	}
	if err := pgxPool.Ping(ctx); err != nil {
		cancel()
		return nil, errors.Wrapf(entities.ErrInternal, "failed to ping database: %v", err)
	}

	return &Storage{
		pool:     pgxPool,
		cancelFn: cancel,
		once:     sync.Once{},
	}, nil
}

func (s *Storage) Close() {
	s.once.Do(func() {
		s.cancelFn()
		s.pool.Close()
	})
}

func (s *Storage) Store(ctx context.Context, coins []*entities.Coin) error { // положить - метод сохранения в хранилище
	if len(coins) == 0 {
		return errors.Wrap(entities.ErrInvalidParam, "no coins")
	}

	sqlQuery := `
	INSERT INTO crypto.coins (title, price, actual_at)
	VALUES ($1, $2, $3);`

	batch := &pgx.Batch{}

	for _, elem := range coins {
		batch.Queue(sqlQuery, elem.Title(), elem.Price(), elem.ActualAt())
	}

	batchResult := s.pool.SendBatch(ctx, batch)
	defer batchResult.Close()

	for range coins {
		if _, err := batchResult.Exec(); err != nil {
			return errors.Wrapf(entities.ErrInternal, "failed to save to storage: %v", err)
		}
	}

	return nil
}

func (s *Storage) GetCoinsList(ctx context.Context) ([]string, error) { // список имен монет, которые представлены
	sqlQuery := `
	SELECT DISTINCT title
	FROM crypto.coins 
	ORDER BY title ASC;`

	rows, err := s.pool.Query(ctx, sqlQuery)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to query storage: %v", err)
	}
	defer rows.Close()

	titleCoins := make([]string, 0, 100)

	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, errors.Wrapf(entities.ErrInternal, "failed to scan rows: %v", err)
		}
		titleCoins = append(titleCoins, title)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to rows iteration: %v", err)
	}

	return titleCoins, nil
}

func (s *Storage) GetActualCoins(ctx context.Context, titles []string) ([]*entities.Coin, error) { // получение последних монет по title-ам
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "title cannot be empty")
	}

	sqlQuery := `
	SELECT DISTINCT ON(title) title, price, actual_at
	FROM crypto.coins
	WHERE title = ANY($1)
	ORDER BY title, actual_at DESC;`

	rows, err := s.pool.Query(ctx, sqlQuery, titles)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to query storage: %v", err)
	}
	defer rows.Close()

	sliceCoins := make([]*entities.Coin, 0, len(titles))

	for rows.Next() {
		var coinModel CryptoModel
		if err := rows.Scan(&coinModel.Title, &coinModel.Price, &coinModel.ActualAT); err != nil {
			return nil, errors.Wrapf(entities.ErrInternal, "failed to scan rows: %v", err)
		}
		coin, err := entities.NewCoin(coinModel.Title, coinModel.Price, coinModel.ActualAT)
		if err != nil {
			return nil, errors.Wrapf(entities.ErrInternal, "failed to create coin model: %v", err)
		}
		sliceCoins = append(sliceCoins, coin)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to rows iteration: %v", err)
	}

	return sliceCoins, nil
}

func (s *Storage) GetAggregateCoins(ctx context.Context, titles []string, aggType string) ([]*entities.Coin, error) { // Агрегированный запрос над монетами
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "title cannot be empty")
	}
	aggTypeUpper := strings.ToUpper(aggType)
	if aggTypeUpper != "MIN" && aggTypeUpper != "MAX" && aggTypeUpper != "AVG" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "incorrect value for aggregated rates")
	}

	sqlQuery := fmt.Sprintf(`
	SELECT title, %s(price) AS result_price
	FROM crypto.coins
	WHERE title = ANY($1)
	GROUP BY title
	ORDER BY result_price ASC;`, aggTypeUpper)

	rows, err := s.pool.Query(ctx, sqlQuery, titles)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to query storage: %v", err)
	}
	defer rows.Close()

	sliceCoins := make([]*entities.Coin, 0, len(titles))
	actuaAt := time.Now()

	for rows.Next() {
		var coinModel CryptoModel
		if err := rows.Scan(&coinModel.Title, &coinModel.Price); err != nil {
			return nil, errors.Wrapf(entities.ErrInternal, "failed to scan rows: %v", err)
		}
		coin, err := entities.NewCoin(coinModel.Title, coinModel.Price, actuaAt)
		if err != nil {
			return nil, errors.Wrapf(entities.ErrInternal, "failed to create coin model: %v", err)
		}
		sliceCoins = append(sliceCoins, coin)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to rows iteration: %v", err)
	}

	return sliceCoins, nil
}
