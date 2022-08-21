package repository_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/interfaceadapters/repository"
	usecaserepository "github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase/repository"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testFavoriteRepositorySuite struct {
	suite.Suite
}

func TestFavoriteRepositorySuite(t *testing.T) {
	suite.Run(t, &testFavoriteRepositorySuite{})
}

func (s *testFavoriteRepositorySuite) TestAddFavoriteMovie() {
	type testInput struct {
		args  usecaserepository.AddFavoriteMovieParams
		mocks func(mock sqlmock.Sqlmock)
	}

	type testOutput struct {
		err error
	}

	cases := []struct {
		name     string
		input    testInput
		expected testOutput
	}{
		{
			name: "returns_nil_when_add_favorite_movie_record_successfully",
			input: testInput{
				args: usecaserepository.AddFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mocks: func(mock sqlmock.Sqlmock) {
					mock.
						ExpectExec(regexp.QuoteMeta("INSERT INTO favorites(user_id, movie_id) VALUES (?,?)")).
						WithArgs(1, 10).
						WillReturnResult(sqlmock.NewResult(1, 1))

				},
			},
			expected: testOutput{
				err: nil,
			},
		},
		{
			name: "returns_error_when_add_favorite_movie_fail",
			input: testInput{
				args: usecaserepository.AddFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mocks: func(mock sqlmock.Sqlmock) {
					mock.
						ExpectExec(regexp.QuoteMeta("INSERT INTO favorites(user_id, movie_id) VALUES (?,?)")).
						WithArgs(1, 10).
						WillReturnError(fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				err: fmt.Errorf("ExecContext: %w", fmt.Errorf("dummy error")),
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			manager, clean := initMockConnManager(t, c.input.mocks)
			defer clean()

			favoriteRepository := repository.NewFavoriteRepository(manager)

			ctx := context.Background()
			err := favoriteRepository.AddFavoriteMovie(ctx, c.input.args)
			assert.Equal(t, c.expected.err, err)
		})
	}
}

func (s *testFavoriteRepositorySuite) TestCheckIsFavoriteMovie() {
	type testInput struct {
		args  usecaserepository.CheckIsFavoriteMovieParams
		mocks func(mock sqlmock.Sqlmock)
	}

	type testOutput struct {
		isFavoriteMovie bool
		err             error
	}

	cases := []struct {
		name     string
		input    testInput
		expected testOutput
	}{
		{
			name: "returns_true_when_there_is_exist_record_with_correspodinng_user_movie_id",
			input: testInput{
				args: usecaserepository.CheckIsFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mocks: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows(favoritesTableRows)
					rows.AddRow(1, 10, time.Now(), time.Now())

					mock.
						ExpectQuery(regexp.QuoteMeta("SELECT user_id, movie_id, created_at, updated_at FROM favorites WHERE user_id = ? AND movie_id = ?")).
						WithArgs(1, 10).
						WillReturnRows(rows)
				},
			},
			expected: testOutput{
				err:             nil,
				isFavoriteMovie: true,
			},
		},
		{
			name: "returns_false_when_there_is_no_exist_record_corresponding_user_movie_id",
			input: testInput{
				args: usecaserepository.CheckIsFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mocks: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows(favoritesTableRows)
					rows.AddRow(1, 10, time.Now(), time.Now())

					mock.
						ExpectQuery(regexp.QuoteMeta("SELECT user_id, movie_id, created_at, updated_at FROM favorites WHERE user_id = ? AND movie_id = ?")).
						WithArgs(1, 10).
						WillReturnRows(&sqlmock.Rows{})
				},
			},
			expected: testOutput{
				err:             nil,
				isFavoriteMovie: false,
			},
		},
		{
			name: "returns_error_when_query_failed",
			input: testInput{
				args: usecaserepository.CheckIsFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mocks: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows(favoritesTableRows)
					rows.AddRow(1, 10, time.Now(), time.Now())

					mock.
						ExpectQuery(regexp.QuoteMeta("SELECT user_id, movie_id, created_at, updated_at FROM favorites WHERE user_id = ? AND movie_id = ?")).
						WithArgs(1, 10).
						WillReturnError(fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				err:             fmt.Errorf("QueryRowxContext: %w", fmt.Errorf("dummy error")),
				isFavoriteMovie: false,
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			manager, clean := initMockConnManager(t, c.input.mocks)
			defer clean()

			favoriteRepository := repository.NewFavoriteRepository(manager)

			ctx := context.Background()
			res, err := favoriteRepository.CheckIsFavoriteMovie(ctx, c.input.args)
			assert.Equal(t, c.expected.isFavoriteMovie, res)
			assert.Equal(t, c.expected.err, err)
		})
	}
}

func (s *testFavoriteRepositorySuite) TestFindFavoriteMoviesByUserID() {
	type testInput struct {
		userID uint64
		mocks  func(mock sqlmock.Sqlmock)
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
			name: "returns_favorite_movies",
			input: testInput{
				userID: 1,
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
					rows.AddRow(
						2,
						"arcu. Vivamus sit amet risus. Donec egestas. Aliquam",
						"Belgium",
						"egestas, urna justo faucibus lectus, a sollicitudin orci sem eget massa. Suspendisse eleifend. Cras sed leo. Cras vehicula aliquet libero. Integer in magna. Phasellus dolor elit, pellentesque a, facilisis non, bibendum sed, est. Nunc laoreet lectus quis massa. Mauris vestibulum, neque sed dictum eleifend, nunc risus varius orci, in consequat enim diam vel arcu. Curabitur ut odio vel est tempor bibendum. Donec felis orci, adipiscing non, luctus sit amet, faucibus ut, nulla.",
						nil,
						nil,
						false,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						100000,
						1000000,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"))
					mock.
						ExpectQuery(regexp.QuoteMeta(`SELECT movies.id, movies.original_title, movies.original_language,
						movies.overview, movies.poster_path, movies.backdrop_path,
						movies.adult, movies.release_date, movies.budget, movies.revenue, movies.created_at, movies.updated_at
						FROM movies
						INNER JOIN favorites
						ON movies.id = favorites.movie_id
						WHERE favorites.user_id = ?
						ORDER BY movies.id ASC`)).
						WithArgs(1).
						WillReturnRows(rows)
				},
			},
			expected: testOutput{
				err: nil,
				movies: []*entity.Movie{
					{
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
					{
						ID:               2,
						OriginalTitle:    "arcu. Vivamus sit amet risus. Donec egestas. Aliquam",
						OriginalLanguage: "Belgium",
						Overview:         utils.StringPtr("egestas, urna justo faucibus lectus, a sollicitudin orci sem eget massa. Suspendisse eleifend. Cras sed leo. Cras vehicula aliquet libero. Integer in magna. Phasellus dolor elit, pellentesque a, facilisis non, bibendum sed, est. Nunc laoreet lectus quis massa. Mauris vestibulum, neque sed dictum eleifend, nunc risus varius orci, in consequat enim diam vel arcu. Curabitur ut odio vel est tempor bibendum. Donec felis orci, adipiscing non, luctus sit amet, faucibus ut, nulla."),
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
			name: "returns_error_when_query_failed",
			input: testInput{
				userID: 1,
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
					rows.AddRow(
						2,
						"arcu. Vivamus sit amet risus. Donec egestas. Aliquam",
						"Belgium",
						"egestas, urna justo faucibus lectus, a sollicitudin orci sem eget massa. Suspendisse eleifend. Cras sed leo. Cras vehicula aliquet libero. Integer in magna. Phasellus dolor elit, pellentesque a, facilisis non, bibendum sed, est. Nunc laoreet lectus quis massa. Mauris vestibulum, neque sed dictum eleifend, nunc risus varius orci, in consequat enim diam vel arcu. Curabitur ut odio vel est tempor bibendum. Donec felis orci, adipiscing non, luctus sit amet, faucibus ut, nulla.",
						nil,
						nil,
						false,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						100000,
						1000000,
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"))
					mock.
						ExpectQuery(regexp.QuoteMeta(`SELECT movies.id, movies.original_title, movies.original_language,
						movies.overview, movies.poster_path, movies.backdrop_path,
						movies.adult, movies.release_date, movies.budget, movies.revenue, movies.created_at, movies.updated_at
						FROM movies
						INNER JOIN favorites
						ON movies.id = favorites.movie_id
						WHERE favorites.user_id = ?
						ORDER BY movies.id ASC`)).
						WithArgs(1).
						WillReturnError(fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				err:    fmt.Errorf("QueryxContext: %w", fmt.Errorf("dummy error")),
				movies: nil,
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			manager, clean := initMockConnManager(t, c.input.mocks)
			defer clean()

			favoriteRepository := repository.NewFavoriteRepository(manager)

			ctx := context.Background()
			res, err := favoriteRepository.FindFavoriteMoviesByUserID(ctx, c.input.userID)
			assert.Equal(t, c.expected.movies, res)
			assert.Equal(t, c.expected.err, err)
		})
	}
}
