package repository_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/interfaceadapters/repository"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testMovieRepositorySuite struct {
	suite.Suite
}

func TestMovieRepositorySuite(t *testing.T) {
	suite.Run(t, &testMovieRepositorySuite{})
}

func (s *testMovieRepositorySuite) TestFindByID() {
	type testInput struct {
		movieID uint64
		mocks   func(mock sqlmock.Sqlmock)
	}

	type testOutput struct {
		movie *entity.Movie
		err   error
	}

	cases := []struct {
		name     string
		input    testInput
		expected testOutput
	}{
		{
			name: "returns_movie_when_exist_record",
			input: testInput{
				movieID: 1,
				mocks: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows(moviesTableRows)
					rows.AddRow(
						1,
						"accumsan sed, facilisis vitae,",
						"Nigeria",
						"risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras",
						nil,
						nil,
						false,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						100000,
						1000000,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"))
					mock.
						ExpectQuery(regexp.QuoteMeta(`SELECT id, original_title, original_language, overview, poster_path, backdrop_path,
					adult, release_date, budget, revenue, created_at, updated_at FROM movies WHERE id = ?`)).
						WillReturnRows(rows)
				},
			},
			expected: testOutput{
				err: nil,
				movie: &entity.Movie{
					ID:               1,
					OriginalTitle:    "accumsan sed, facilisis vitae,",
					OriginalLanguage: "Nigeria",
					Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras"),
					Adult:            false,
					ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
					Revenue:          utils.Int64Ptr(1000000),
					Budget:           utils.Uint64Ptr(100000),
					CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
					UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
				},
			},
		},
		{
			name: "returns_error_when_query_failed",
			input: testInput{
				movieID: 1,
				mocks: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows(moviesTableRows)
					rows.AddRow(
						1,
						"accumsan sed, facilisis vitae,",
						"Nigeria",
						"risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras",
						nil,
						nil,
						false,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						100000,
						1000000,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"))
					mock.
						ExpectQuery(regexp.QuoteMeta(`SELECT id, original_title, original_language, overview, poster_path, backdrop_path,
					adult, release_date, budget, revenue, created_at, updated_at FROM movies WHERE id = ?`)).
						WillReturnError(fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				err: fmt.Errorf("QueryRowxContext: %w", fmt.Errorf("dummy error")),
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			manager, clean := initMockConnManager(t, c.input.mocks)
			defer clean()

			movieRepository := repository.NewMovieRepository(manager)

			ctx := context.Background()
			res, err := movieRepository.FindByID(ctx, c.input.movieID)
			assert.Equal(t, c.expected.movie, res)
			assert.Equal(t, c.expected.err, err)
		})
	}
}

func (s *testMovieRepositorySuite) TestFindByKeyword() {
	type testInput struct {
		keyword string
		mocks   func(mock sqlmock.Sqlmock)
	}

	type testOutput struct {
		movies []*entity.Movie
		err    error
	}

	cases := []struct {
		name     string
		input    testInput
		expected testOutput
	}{
		{
			name: "returns_movies_match_keyword",
			input: testInput{
				keyword: "test",
				mocks: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows(moviesTableRows)
					rows.AddRow(
						1,
						"test sed, facilisis vitae,",
						"Nigeria",
						"risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras",
						nil,
						nil,
						false,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						100000,
						1000000,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"))
					rows.AddRow(
						2,
						"sed, facilisis vitae,",
						"test",
						"risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras",
						nil,
						nil,
						false,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						100000,
						1000000,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"))
					mock.
						ExpectQuery(regexp.QuoteMeta(`SELECT id, original_title, original_language, overview, poster_path, backdrop_path,
					adult, release_date, budget, revenue, created_at, updated_at
					FROM movies
					WHERE MATCH (original_title, overview, original_language) AGAINST ('test*' IN BOOLEAN MODE)
					ORDER BY id ASC`)).
						WillReturnRows(rows)
				},
			},
			expected: testOutput{
				movies: []*entity.Movie{
					{
						ID:               1,
						OriginalTitle:    "test sed, facilisis vitae,",
						OriginalLanguage: "Nigeria",
						Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras"),
						Adult:            false,
						ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
						Revenue:          utils.Int64Ptr(1000000),
						Budget:           utils.Uint64Ptr(100000),
						CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
					},
					{
						ID:               2,
						OriginalTitle:    "sed, facilisis vitae,",
						OriginalLanguage: "test",
						Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras"),
						Adult:            false,
						ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
						Revenue:          utils.Int64Ptr(1000000),
						Budget:           utils.Uint64Ptr(100000),
						CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
					},
				},
			},
		},
		{
			name: "returns_errors_when_query_failed",
			input: testInput{
				keyword: "test",
				mocks: func(mock sqlmock.Sqlmock) {
					mock.
						ExpectQuery(regexp.QuoteMeta(`SELECT id, original_title, original_language, overview, poster_path, backdrop_path,
					adult, release_date, budget, revenue, created_at, updated_at
					FROM movies
					WHERE MATCH (original_title, overview, original_language) AGAINST ('test*' IN BOOLEAN MODE)
					ORDER BY id ASC`)).
						WillReturnError(fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				movies: nil,
				err:    fmt.Errorf("QueryxContext: %w", fmt.Errorf("dummy error")),
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			manager, clean := initMockConnManager(t, c.input.mocks)
			defer clean()

			movieRepository := repository.NewMovieRepository(manager)

			ctx := context.Background()
			res, err := movieRepository.FindByKeyword(ctx, c.input.keyword)
			assert.Equal(t, c.expected.movies, res)
			assert.Equal(t, c.expected.err, err)
		})
	}
}

func (s *testMovieRepositorySuite) TestFindPopularMovies() {
	type testInput struct {
		limit uint64
		mocks func(mock sqlmock.Sqlmock)
	}

	type testOutput struct {
		movies []*entity.Movie
		err    error
	}

	cases := []struct {
		name     string
		input    testInput
		expected testOutput
	}{
		{
			name: "returns_popular_movies",
			input: testInput{
				limit: 10,
				mocks: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows(moviesTableRows)
					rows.AddRow(
						1,
						"test sed, facilisis vitae,",
						"Nigeria",
						"risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras",
						nil,
						nil,
						false,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						100000,
						1000000,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"))
					rows.AddRow(
						2,
						"sed, facilisis vitae,",
						"test",
						"risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras",
						nil,
						nil,
						false,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						100000,
						1000000,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"))
					mock.
						ExpectQuery(regexp.QuoteMeta(`SELECT id, original_title, original_language, overview, poster_path, backdrop_path,
					adult, release_date, budget, revenue, created_at, updated_at, IFNULL(favorite_numbers.favorite_number, 0) as favorite_number
					FROM movies
					LEFT JOIN (SELECT movie_id, count(*) AS favorite_number FROM favorites GROUP BY movie_id) AS favorite_numbers
					ON movies.id = favorite_numbers.movie_id
					ORDER BY favorite_number DESC
					LIMIT 10`)).
						WillReturnRows(rows)
				},
			},
			expected: testOutput{
				movies: []*entity.Movie{
					{
						ID:               1,
						OriginalTitle:    "test sed, facilisis vitae,",
						OriginalLanguage: "Nigeria",
						Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras"),
						Adult:            false,
						ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
						Revenue:          utils.Int64Ptr(1000000),
						Budget:           utils.Uint64Ptr(100000),
						CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
					},
					{
						ID:               2,
						OriginalTitle:    "sed, facilisis vitae,",
						OriginalLanguage: "test",
						Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras"),
						Adult:            false,
						ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
						Revenue:          utils.Int64Ptr(1000000),
						Budget:           utils.Uint64Ptr(100000),
						CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
					},
				},
			},
		},
		{
			name: "returns_errors_when_query_failed",
			input: testInput{
				limit: 10,
				mocks: func(mock sqlmock.Sqlmock) {
					mock.
						ExpectQuery(regexp.QuoteMeta(`SELECT id, original_title, original_language, overview, poster_path, backdrop_path,
					adult, release_date, budget, revenue, created_at, updated_at, IFNULL(favorite_numbers.favorite_number, 0) as favorite_number
					FROM movies
					LEFT JOIN (SELECT movie_id, count(*) AS favorite_number FROM favorites GROUP BY movie_id) AS favorite_numbers
					ON movies.id = favorite_numbers.movie_id
					ORDER BY favorite_number DESC
					LIMIT 10`)).
						WillReturnError(fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				movies: nil,
				err:    fmt.Errorf("QueryxContext: %w", fmt.Errorf("dummy error")),
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			manager, clean := initMockConnManager(t, c.input.mocks)
			defer clean()

			movieRepository := repository.NewMovieRepository(manager)

			ctx := context.Background()
			res, err := movieRepository.FindPopularMovies(ctx, uint(c.input.limit))
			assert.Equal(t, c.expected.movies, res)
			assert.Equal(t, c.expected.err, err)
		})
	}
}
