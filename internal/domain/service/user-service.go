package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"playoo/internal/domain/entity"
	"playoo/internal/domain/repository"
	"strings"
	"math"
)

type UserService interface {
	GetById(c context.Context, id string) (*entity.User, error)
	GetByAmount(c context.Context, amount int) ([]entity.User,error)
	UpdateEvents(c context.Context, id,eventid string) error
	UploadAvatar(c context.Context, id string,picture *multipart.FileHeader) error
	DeleteAvatar(c context.Context,id string) error
	RecordDiscord(c context.Context, id,ds string) error
	EditRating(c context.Context,id string,stars int) error
	Follow(c context.Context,id,login string) error
	Unfollow(c context.Context,id,login string) error
}

type userService struct{
	UserRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService{
	return &userService{
		UserRepository: userRepository,
	}
}	

func (us userService) GetById(c context.Context, id string) (*entity.User, error){
	return us.UserRepository.FindById(c,id)
}

func (us userService) GetByAmount(c context.Context, amount int) ([]entity.User,error){
	return us.UserRepository.FindByAmount(c,amount)
}


func (ur userService) UpdateEvents(c context.Context, id,eventid string) error{
	user,err:=ur.UserRepository.FindById(c,id)
	if err!=nil{
		return err
	}
	updateEvents:=make([]string,0,len(user.Events))
	for _, e := range user.Events {
		if e != eventid{
			updateEvents = append(updateEvents, e)
		}
	}
	user.Events = updateEvents
	if err:=ur.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	return nil
}

func (ur userService) UploadAvatar(c context.Context, id string, picture *multipart.FileHeader) error{
	user,err:=ur.GetById(c,id)
	if err!=nil{
		return err
	}
	uploadDir := "../../files/avatars"
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		return err
	}
	fileName := fmt.Sprintf("%s%s", user.Id,filepath.Ext(picture.Filename))
    filepath := filepath.Join(uploadDir, fileName)
    dst, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer dst.Close()
	src,err:=picture.Open()
	if err!=nil{
		return err
	}
	defer src.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	var(
		host = os.Getenv("HOST")
		port = os.Getenv("PORT")
	)
	fileURL:=fmt.Sprintf("http://%s:%s/files/avatars/%s",host,port,fileName)

	user.Avatar = fileURL
	if err:=ur.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	return nil
}

func(ur userService) DeleteAvatar(c context.Context,id string) error{
	user,err:=ur.GetById(c,id)
	if err!=nil{
		return err
	}
	if user.Avatar == "" {
		return errors.New("user does not have an avatar")
	}
	var(
		host = os.Getenv("HOST")
		port = os.Getenv("PORT")
	)
	file:=strings.TrimPrefix(user.Avatar,fmt.Sprintf("http://%s:%s/",host,port))
	if err:=os.Remove(fmt.Sprintf("../../%s",file));err!=nil{
		return err
	}
	user.Avatar = ""
	if err:=ur.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	return nil
}

func(ur userService) RecordDiscord(c context.Context, id,ds string) error{
	user,err:=ur.GetById(c,id)
	if err!=nil{
		return err
	}
	user.Discord=ds
	if err:=ur.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	return nil
}

func(ur userService) EditRating(c context.Context,id string,stars int) error{
	user,err:=ur.GetById(c,id)
	if err!=nil{
		return err
	}
	user.NumberOfRatings++
	user.TotalRating += stars
	averageRating := float64(user.TotalRating) / float64(user.NumberOfRatings)
	user.Rating = math.Round(averageRating*2)/2
	if err:=ur.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	return nil
}

func (ur userService) Follow(c context.Context,id,login string) error{
	user,err:=ur.UserRepository.FindById(c,id)
	if err!=nil{
		return err
	}
	follow,err:=ur.UserRepository.FindByLogin(c,login)
	if err!=nil{
		return err
	}
	for _,followId:=range user.Followings{
		if followId == follow.Id.String(){
			return errors.New("you alredy follow to this user")
		}
	}
	
	user.Followings = append(user.Followings, follow.Id.String())
	follow.Followers = append(follow.Followers, user.Id.String())

	if err:=ur.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	if err:=ur.UserRepository.Save(c,*follow);err!=nil{
		return err
	}
	return nil
}

func (ur userService) Unfollow(c context.Context,id,login string) error{
	user,err:=ur.UserRepository.FindById(c,id)
	if err!=nil{
		return err
	}
	follow,err:=ur.UserRepository.FindByLogin(c,login)
	if err!=nil{
		return err
	}
	updateFollowings:=make([]string,0,len(user.Followings))
	for _, f := range user.Followings {
		if f != follow.Id.String() {
			updateFollowings = append(updateFollowings, f)
		}
	}
	updateFollowers:=make([]string,0,len(follow.Followers))
	for _, f := range user.Followers {
		if f !=  user.Id.String(){
			updateFollowers = append(updateFollowers, f)
		}
	}
	if len(updateFollowings) != 0{
		user.Followings = updateFollowings
	}else{
		user.Followings = nil
	}
	if len(updateFollowers)!=0{
		follow.Followers = updateFollowers
	}else{
		follow.Followers = nil
	}
	
	user.Followings = append(user.Followings, follow.Id.String())
	follow.Followers = append(follow.Followers, user.Id.String())

	if err:=ur.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	if err:=ur.UserRepository.Save(c,*follow);err!=nil{
		return err
	}
	return nil
}