package usecase_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/samthehai/ml-backend-test-samthehai/config"
	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase/repository"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase/testdata/mock_repository"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/httperrors"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testMovieUsecase struct {
	suite.Suite
}

func TestMovieUsecasesuite(t *testing.T) {
	suite.Run(t, &testMovieUsecase{})
}

func (s *testMovieUsecase) TestGetMovieByID() {
	type testInput struct {
		movieID             uint64
		mockMovieRepository func(*mock_repository.MockMovieRepository)
	}

	type testOutput struct {
		err   error
		movie *entity.Movie
	}

	cases := []struct {
		name     string
		input    testInput
		expected testOutput
	}{
		{
			name: "returns_movie_when_found",
			input: testInput{
				movieID: 1,
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByID(gomock.Any(), uint64(1)).Return(
						&entity.Movie{
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
						}, nil,
					)
				},
			},
			expected: testOutput{
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
			name: "returns_notfound_error_when_not_found",
			input: testInput{
				movieID: 1,
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByID(gomock.Any(), uint64(1)).Return(nil, nil)
				},
			},
			expected: testOutput{
				movie: nil,
				err:   httperrors.NewNotFoundError(fmt.Errorf("movieRepository.FindByID: not found")),
			},
		},
		{
			name: "returns_error_of_FindByID_when_error_happended",
			input: testInput{
				movieID: 1,
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByID(gomock.Any(), uint64(1)).Return(nil, fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				movie: nil,
				err:   httperrors.NewInternalServerError(fmt.Errorf("movieRepository.FindByID: %w", fmt.Errorf("dummy error"))),
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMovieRepository := mock_repository.NewMockMovieRepository(ctrl)
			c.input.mockMovieRepository(mockMovieRepository)

			u := usecase.NewMovieUsecase(config.Config{}, logger.NewApiLogger(&config.Config{}), mockMovieRepository, nil)
			res, err := u.GetMovieByID(context.Background(), c.input.movieID)
			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.movie, res)
		})
	}
}

func (s *testMovieUsecase) TestSearchByKeyword() {
	type testInput struct {
		keyword             string
		mockMovieRepository func(*mock_repository.MockMovieRepository)
	}

	type testOutput struct {
		err    error
		movies []*entity.Movie
	}

	cases := []struct {
		name     string
		input    testInput
		expected testOutput
	}{
		{
			name: "returns_popular_movies_when_keyword_is_empty",
			input: testInput{
				keyword: "",
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindPopularMovies(gomock.Any(), uint(100)).Return(
						[]*entity.Movie{
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
								OriginalTitle:    "accumsan sed, facilisis vitae,2",
								OriginalLanguage: "Nigeria",
								Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras22"),
								Adult:            false,
								ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
								Revenue:          utils.Int64Ptr(1000000),
								Budget:           utils.Uint64Ptr(100000),
								CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
								UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
							},
						}, nil,
					)
				},
			},
			expected: testOutput{
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
						OriginalTitle:    "accumsan sed, facilisis vitae,2",
						OriginalLanguage: "Nigeria",
						Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras22"),
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
			name: "returns_error_of_FindPopularMovies_when_error_happended",
			input: testInput{
				keyword: "",
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindPopularMovies(gomock.Any(), uint(100)).Return(nil, fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				movies: nil,
				err:    httperrors.NewInternalServerError(fmt.Errorf("movieRepository.FindPopularMovies: %w", fmt.Errorf("dummy error"))),
			},
		},
		{
			name: "returns_movies_when_keyword_is_not_empty",
			input: testInput{
				keyword: "test",
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByKeyword(gomock.Any(), "test").Return(
						[]*entity.Movie{
							{
								ID:               1,
								OriginalTitle:    "accumsan tested, facilisis vitae,",
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
								OriginalTitle:    "accumsan sed, test facilisis vitae,2",
								OriginalLanguage: "Nigeria",
								Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras22"),
								Adult:            false,
								ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
								Revenue:          utils.Int64Ptr(1000000),
								Budget:           utils.Uint64Ptr(100000),
								CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
								UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
							},
						}, nil,
					)
				},
			},
			expected: testOutput{
				movies: []*entity.Movie{
					{
						ID:               1,
						OriginalTitle:    "accumsan tested, facilisis vitae,",
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
						OriginalTitle:    "accumsan sed, test facilisis vitae,2",
						OriginalLanguage: "Nigeria",
						Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras22"),
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
			name: "returns_error_of_FindByKeyword_when_error_happended",
			input: testInput{
				keyword: "test",
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByKeyword(gomock.Any(), "test").Return(nil, fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				movies: nil,
				err:    httperrors.NewInternalServerError(fmt.Errorf("movieRepository.FindByKeyword: %w", fmt.Errorf("dummy error"))),
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMovieRepository := mock_repository.NewMockMovieRepository(ctrl)
			c.input.mockMovieRepository(mockMovieRepository)

			u := usecase.NewMovieUsecase(config.Config{}, logger.NewApiLogger(&config.Config{}), mockMovieRepository, nil)
			res, err := u.SearchByKeyword(context.Background(), c.input.keyword)
			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.movies, res)
		})
	}
}

func (s *testMovieUsecase) TestAddFavoriteMovie() {
	type testInput struct {
		args                   repository.AddFavoriteMovieParams
		mockMovieRepository    func(*mock_repository.MockMovieRepository)
		mockFavoriteRepository func(*mock_repository.MockFavoriteRepository)
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
			name: "returns_nil_when_add_favorite_movie_successfully",
			input: testInput{
				args: repository.AddFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByID(gomock.Any(), uint64(10)).Return(
						&entity.Movie{
							ID:               10,
							OriginalTitle:    "accumsan tested, facilisis vitae,",
							OriginalLanguage: "Nigeria",
							Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras"),
							Adult:            false,
							ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
							Revenue:          utils.Int64Ptr(1000000),
							Budget:           utils.Uint64Ptr(100000),
							CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
							UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						}, nil,
					)
				},
				mockFavoriteRepository: func(r *mock_repository.MockFavoriteRepository) {
					r.EXPECT().CheckIsFavoriteMovie(gomock.Any(), repository.CheckIsFavoriteMovieParams{
						UserID:  1,
						MovieID: 10,
					}).Return(false, nil)
					r.EXPECT().AddFavoriteMovie(gomock.Any(), repository.AddFavoriteMovieParams{
						UserID:  1,
						MovieID: 10,
					}).Return(nil)
				},
			},
			expected: testOutput{
				err: nil,
			},
		},
		{
			name: "returns_error_of_FindByID_when_FindByID_returns_error",
			input: testInput{
				args: repository.AddFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByID(gomock.Any(), uint64(10)).Return(nil, fmt.Errorf("dummy error"))
				},
				mockFavoriteRepository: func(mfr *mock_repository.MockFavoriteRepository) {},
			},
			expected: testOutput{
				err: httperrors.NewInternalServerError(fmt.Errorf("movieRepository.FindByID: %w", fmt.Errorf("dummy error"))),
			},
		},
		{
			name: "returns_error_when_not_found_movie",
			input: testInput{
				args: repository.AddFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByID(gomock.Any(), uint64(10)).Return(nil, nil)
				},
				mockFavoriteRepository: func(mfr *mock_repository.MockFavoriteRepository) {},
			},
			expected: testOutput{
				err: httperrors.NewNotFoundError(fmt.Errorf("movieRepository.FindByID: not found")),
			},
		},
		{
			name: "returns_error_of_CheckIsFavoriteMovie_when_CheckIsFavoriteMovie_returns_error",
			input: testInput{
				args: repository.AddFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByID(gomock.Any(), uint64(10)).Return(
						&entity.Movie{
							ID:               10,
							OriginalTitle:    "accumsan tested, facilisis vitae,",
							OriginalLanguage: "Nigeria",
							Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras"),
							Adult:            false,
							ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
							Revenue:          utils.Int64Ptr(1000000),
							Budget:           utils.Uint64Ptr(100000),
							CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
							UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						}, nil,
					)
				},
				mockFavoriteRepository: func(r *mock_repository.MockFavoriteRepository) {
					r.EXPECT().CheckIsFavoriteMovie(gomock.Any(), repository.CheckIsFavoriteMovieParams{
						UserID:  1,
						MovieID: 10,
					}).Return(false, fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				err: httperrors.NewInternalServerError(fmt.Errorf("favoriteRepository.CheckIsFavoriteMovie: %w", fmt.Errorf("dummy error"))),
			},
		},
		{
			name: "returns_error_already_favorite_when_movie_already_is_favor_by_user",
			input: testInput{
				args: repository.AddFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByID(gomock.Any(), uint64(10)).Return(
						&entity.Movie{
							ID:               10,
							OriginalTitle:    "accumsan tested, facilisis vitae,",
							OriginalLanguage: "Nigeria",
							Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras"),
							Adult:            false,
							ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
							Revenue:          utils.Int64Ptr(1000000),
							Budget:           utils.Uint64Ptr(100000),
							CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
							UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						}, nil,
					)
				},
				mockFavoriteRepository: func(r *mock_repository.MockFavoriteRepository) {
					r.EXPECT().CheckIsFavoriteMovie(gomock.Any(), repository.CheckIsFavoriteMovieParams{
						UserID:  1,
						MovieID: 10,
					}).Return(true, nil)
				},
			},
			expected: testOutput{
				err: httperrors.NewRestError(http.StatusBadRequest, "already is favorited", nil),
			},
		},
		{
			name: "returns_error_of_AddFavoriteMovie",
			input: testInput{
				args: repository.AddFavoriteMovieParams{
					UserID:  1,
					MovieID: 10,
				},
				mockMovieRepository: func(r *mock_repository.MockMovieRepository) {
					r.EXPECT().FindByID(gomock.Any(), uint64(10)).Return(
						&entity.Movie{
							ID:               10,
							OriginalTitle:    "accumsan tested, facilisis vitae,",
							OriginalLanguage: "Nigeria",
							Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras"),
							Adult:            false,
							ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
							Revenue:          utils.Int64Ptr(1000000),
							Budget:           utils.Uint64Ptr(100000),
							CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
							UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						}, nil,
					)
				},
				mockFavoriteRepository: func(r *mock_repository.MockFavoriteRepository) {
					r.EXPECT().CheckIsFavoriteMovie(gomock.Any(), repository.CheckIsFavoriteMovieParams{
						UserID:  1,
						MovieID: 10,
					}).Return(false, nil)
					r.EXPECT().AddFavoriteMovie(gomock.Any(), repository.AddFavoriteMovieParams{
						UserID:  1,
						MovieID: 10,
					}).Return(fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				err: httperrors.NewRestError(http.StatusInternalServerError,
					fmt.Errorf("favoriteRepository.AddFavoriteMovie: %w", fmt.Errorf("dummy error")).Error(), nil),
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMovieRepository := mock_repository.NewMockMovieRepository(ctrl)
			mockFavoriteRepository := mock_repository.NewMockFavoriteRepository(ctrl)
			c.input.mockMovieRepository(mockMovieRepository)
			c.input.mockFavoriteRepository(mockFavoriteRepository)

			u := usecase.NewMovieUsecase(config.Config{}, logger.NewApiLogger(&config.Config{}), mockMovieRepository, mockFavoriteRepository)
			err := u.AddFavoriteMovie(context.Background(), usecase.AddFavoriteMovieParams(c.input.args))
			assert.Equal(t, c.expected.err, err)
		})
	}
}

func (s *testMovieUsecase) TestListFavoriteMoviesByUserID() {
	type testInput struct {
		userID                 uint64
		mockFavoriteRepository func(*mock_repository.MockFavoriteRepository)
	}

	type testOutput struct {
		err    error
		movies []*entity.Movie
	}

	cases := []struct {
		name     string
		input    testInput
		expected testOutput
	}{
		{
			name: "returns_favorite_movies_of_user",
			input: testInput{
				userID: 1,
				mockFavoriteRepository: func(r *mock_repository.MockFavoriteRepository) {
					r.EXPECT().FindFavoriteMoviesByUserID(gomock.Any(), uint64(1)).
						Return(
							[]*entity.Movie{
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
									OriginalTitle:    "accumsan sed, facilisis vitae,2",
									OriginalLanguage: "Nigeria",
									Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras22"),
									Adult:            false,
									ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
									Revenue:          utils.Int64Ptr(1000000),
									Budget:           utils.Uint64Ptr(100000),
									CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
									UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
								},
							}, nil,
						)
				},
			},
			expected: testOutput{
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
						OriginalTitle:    "accumsan sed, facilisis vitae,2",
						OriginalLanguage: "Nigeria",
						Overview:         utils.StringPtr("risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras22"),
						Adult:            false,
						ReleaseDate:      utils.TimePtr(utils.MustRFC3339Time("2022-08-20T22:00:00+00:00")),
						Revenue:          utils.Int64Ptr(1000000),
						Budget:           utils.Uint64Ptr(100000),
						CreatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
						UpdatedAt:        utils.MustRFC3339Time("2022-08-20T22:00:00+00:00"),
					},
				},
				err: nil,
			},
		},
		{
			name: "returns_error_of_FindFavoriteMoviesByUserID_when_it_happended",
			input: testInput{
				userID: 1,
				mockFavoriteRepository: func(r *mock_repository.MockFavoriteRepository) {
					r.EXPECT().FindFavoriteMoviesByUserID(gomock.Any(), uint64(1)).Return(nil, fmt.Errorf("dummy error"))
				},
			},
			expected: testOutput{
				movies: nil,
				err:    httperrors.NewInternalServerError(fmt.Errorf("favoriteRepository.FindFavoriteMoviesByUserID: %w", fmt.Errorf("dummy error"))),
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFavoriteRepository := mock_repository.NewMockFavoriteRepository(ctrl)
			c.input.mockFavoriteRepository(mockFavoriteRepository)

			u := usecase.NewMovieUsecase(config.Config{}, logger.NewApiLogger(&config.Config{}), nil, mockFavoriteRepository)
			res, err := u.ListFavoriteMoviesByUserID(context.Background(), c.input.userID)
			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.movies, res)
		})
	}
}
