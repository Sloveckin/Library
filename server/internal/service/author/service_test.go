package service_author

import (
	"errors"
	"reflect"
	"server/internal/model"
	"server/internal/service/author/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
)

const dbErrorMessage = "db error"

func TestAuthorServiceImplCreate(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Author
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				name: "Mark",
			},
			want:    &model.Author{Name: "Mark"},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				name: "Elon",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := mocks.NewAuthorRepository(t)

			if tt.name == "success" {

				mockRepo.
					On("ExistsByName", tt.args.name).
					Return(false, nil)

				mockRepo.
					On("Create", mock.Anything).
					Return(&model.Author{Name: tt.args.name}, nil)

			} else {

				mockRepo.
					On("ExistsByName", tt.args.name).
					Return(true, nil)

			}

			s := AuthorServiceImpl{
				authorRepository: mockRepo,
			}

			got, err := s.Create(tt.args.name)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorServiceImplDelete(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				id: "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := mocks.NewAuthorRepository(t)

			mockRepo.
				On("ExistsById", tt.args.id).
				Return(true, nil)

			mockRepo.On("Delete", tt.args.id).Return(nil)

			s := &AuthorServiceImpl{
				authorRepository: mockRepo,
			}

			if err := s.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthorServiceImplExistsById(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				id: "1",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := mocks.NewAuthorRepository(t)

			mockRepo.
				On("ExistsById", tt.args.id).
				Return(true, nil)

			s := &AuthorServiceImpl{
				authorRepository: mockRepo,
			}

			got, err := s.ExistsById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExistsById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExistsById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorServiceImplExistsByName(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				id: "Mark",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := mocks.NewAuthorRepository(t)

			mockRepo.
				On("ExistsByName", tt.args.id).
				Return(true, nil)

			s := &AuthorServiceImpl{
				authorRepository: mockRepo,
			}

			got, err := s.ExistsByName(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExistsById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExistsById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorServiceImplGet(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Author
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				id: "1",
			},
			want:    &model.Author{Id: "1", Name: "Mark"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := mocks.NewAuthorRepository(t)
			mockRepo.On("Get", tt.args.id).Return(tt.want, nil)

			s := &AuthorServiceImpl{
				authorRepository: mockRepo,
			}

			got, err := s.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type updateArgs struct {
	id   string
	name string
}

type updateTestCase struct {
	name        string
	args        updateArgs
	setupMock   func(repo *mocks.AuthorRepository)
	want        *model.Author
	wantErr     bool
	expectedErr error
}

func assertUpdateResult(t *testing.T, tt updateTestCase, got *model.Author, err error) {
	t.Helper()

	if tt.wantErr {
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if tt.expectedErr != nil && err != tt.expectedErr {
			t.Errorf("expected err %v, got %v", tt.expectedErr, err)
		}
		return
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(got, tt.want) {
		t.Errorf("Update() got = %v, want %v", got, tt.want)
	}
}

func TestAuthorServiceImplUpdate(t *testing.T) {
	tests := []updateTestCase{
		{
			name: "success update",
			args: updateArgs{
				id:   "1",
				name: "NewName",
			},
			setupMock: func(repo *mocks.AuthorRepository) {
				repo.On("ExistsById", "1").Return(true, nil)
				repo.On("ExistsByName", "NewName").Return(false, nil)
				repo.On("Update", "1", "NewName").
					Return(&model.Author{Id: "1", Name: "NewName"}, nil)
			},
			want: &model.Author{
				Id:   "1",
				Name: "NewName",
			},
			wantErr: false,
		},
		{
			name: "author not exists by id",
			args: updateArgs{
				id:   "1",
				name: "NewName",
			},
			setupMock: func(repo *mocks.AuthorRepository) {
				repo.On("ExistsById", "1").Return(false, nil)
			},
			wantErr:     true,
			expectedErr: ErrAuthorNotExists,
		},
		{
			name: "error from ExistsById",
			args: updateArgs{
				id:   "1",
				name: "NewName",
			},
			setupMock: func(repo *mocks.AuthorRepository) {
				repo.On("ExistsById", "1").Return(false, errors.New(dbErrorMessage))
			},
			wantErr: true,
		},
		{
			name: "error from ExistsByName",
			args: updateArgs{
				id:   "1",
				name: "NewName",
			},
			setupMock: func(repo *mocks.AuthorRepository) {
				repo.On("ExistsById", "1").Return(true, nil)
				repo.On("ExistsByName", "NewName").Return(false, errors.New(dbErrorMessage))
			},
			wantErr: true,
		},
		{
			name: "name already exists for another author",
			args: updateArgs{
				id:   "1",
				name: "ExistingName",
			},
			setupMock: func(repo *mocks.AuthorRepository) {
				repo.On("ExistsById", "1").Return(true, nil)
				repo.On("ExistsByName", "ExistingName").Return(true, nil)
				repo.On("GetByName", "ExistingName").
					Return(&model.Author{Id: "2", Name: "ExistingName"}, nil)
			},
			wantErr:     true,
			expectedErr: ErrAuthorExists,
		},
		{
			name: "error from GetByName",
			args: updateArgs{
				id:   "1",
				name: "ExistingName",
			},
			setupMock: func(repo *mocks.AuthorRepository) {
				repo.On("ExistsById", "1").Return(true, nil)
				repo.On("ExistsByName", "ExistingName").Return(true, nil)
				repo.On("GetByName", "ExistingName").
					Return(nil, errors.New(dbErrorMessage))
			},
			wantErr: true,
		},
		{
			name: "same author updates same name",
			args: updateArgs{
				id:   "1",
				name: "SameName",
			},
			setupMock: func(repo *mocks.AuthorRepository) {
				repo.On("ExistsById", "1").Return(true, nil)
				repo.On("ExistsByName", "SameName").Return(true, nil)
				repo.On("GetByName", "SameName").
					Return(&model.Author{Id: "1", Name: "SameName"}, nil)

				repo.On("Update", "1", "SameName").
					Return(&model.Author{Id: "1", Name: "SameName"}, nil)
			},
			want: &model.Author{
				Id:   "1",
				Name: "SameName",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := mocks.NewAuthorRepository(t)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			s := &AuthorServiceImpl{
				authorRepository: mockRepo,
			}

			got, err := s.Update(tt.args.id, tt.args.name)
			assertUpdateResult(t, tt, got, err)
		})
	}
}
