package services

import (
	"context"
	"errors"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/utils"
	"github.com/habibmrizki/BE-EventHub/pkg"
)

type UserService struct {
	repo       *repositories.UserRepository
	hashConfig *pkg.HashConfig
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo:       repo,
		hashConfig: pkg.NewHashConfig(),
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID int) (dto.UserProfileResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return dto.UserProfileResponse{}, err
	}

	return dto.UserProfileResponse{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		Role:       user.Role,
		Location:   user.Location,
		AvatarURL:  user.AvatarURL,
		Bio:        user.Bio,
		JoinedDate: user.CreatedAt.Format("January 2006"),
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

// func (s *UserService) SaveAvatarFile(fileHeader *multipart.FileHeader, userID int) (string, error) {
// 	//  Validasi Ukuran File (Maksimal 2MB = 2 * 1024 * 1024 bytes)
// 	const maxFileSize = 2 * 1024 * 1024
// 	if fileHeader.Size > maxFileSize {
// 		return "", errors.New("ukuran gambar terlalu besar, maksimal 2MB")
// 	}

// 	//  Validasi Ekstensi Gambar
// 	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
// 	allowedExtensions := map[string]bool{
// 		".jpg":  true,
// 		".jpeg": true,
// 		".png":  true,
// 		".webp": true,
// 	}
// 	if !allowedExtensions[ext] {
// 		return "", errors.New("format file tidak didukung, gunakan JPG, PNG, atau WebP")
// 	}

// 	//  Buat direktori upload jika belum ada
// 	uploadDir := "./uploads/avatars"
// 	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
// 		return "", errors.New("gagal membuat direktori upload")
// 	}

// 	//  Generate nama file unik
// 	filename := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().UnixNano(), ext)
// 	dstPath := filepath.Join(uploadDir, filename)

// 	src, err := fileHeader.Open()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer src.Close()

// 	dst, err := os.Create(dstPath)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer dst.Close()

// 	if _, err = io.Copy(dst, src); err != nil {
// 		return "", err
// 	}

// 	// Kembalikan URL path yang bisa diakses
// 	relativeURL := fmt.Sprintf("/uploads/avatars/%s", filename)
// 	return relativeURL, nil
// }

func (s *UserService) UpdateProfile(ctx context.Context, userID int, req dto.UpdateProfileRequest) (dto.UserProfileResponse, error) {
	if req.FullName != nil && len(*req.FullName) < 3 {
		return dto.UserProfileResponse{}, errors.New("nama lengkap minimal 3 karakter")
	}

	user, err := s.repo.UpdateProfile(ctx, userID, req)
	if err != nil {
		return dto.UserProfileResponse{}, err
	}

	return dto.UserProfileResponse{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		Role:       user.Role,
		Location:   user.Location,
		AvatarURL:  user.AvatarURL,
		Bio:        user.Bio,
		JoinedDate: user.CreatedAt.Format("January 2006"),
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID int, req dto.ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("konfirmasi password baru tidak cocok")
	}

	if req.OldPassword == req.NewPassword {
		return errors.New("password baru tidak boleh sama dengan password lama")
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return err
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verifikasi password lama dengan pkg/hash.go
	isMatch, err := s.hashConfig.CompareHashAndPassword(req.OldPassword, user.Password)
	if err != nil || !isMatch {
		return errors.New("password lama tidak sesuai")
	}

	//  Hash password baru dengan pkg/hash.go
	hashedPassword, err := s.hashConfig.GenHash(req.NewPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, hashedPassword)
}
