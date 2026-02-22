package service_author

import (
	"Library/internal/model"
	"Library/internal/service/author/mocks"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestAuthorServiceImpl_Create(t *testing.T) {
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

func TestAuthorServiceImpl_Delete(t *testing.T) {
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

func TestAuthorServiceImpl_ExistsById(t *testing.T) {
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

func TestAuthorServiceImpl_ExistsByName(t *testing.T) {
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

func TestAuthorServiceImpl_Get(t *testing.T) {
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
