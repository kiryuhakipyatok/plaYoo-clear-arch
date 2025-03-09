package service

import (
	"playoo/internal/domain/repository"
	"context"
	"playoo/internal/domain/entity"
)

type GameService interface {
	AddGameToUser(c context.Context, name,id string) error
	GetByName(c context.Context, name string) (*entity.Game,error)
	GetByAmount(c context.Context, amount int) ([]entity.Game,error)
}

type gameService struct {
	GameRepository 			repository.GameRepository
	UserRepository 			repository.UserRepository
	//TransactionRepository 	repository.TransactionRepository
}

func NewGameService(gameRepository repository.GameRepository,userRepository repository.UserRepository) GameService{
	return &gameService{
		GameRepository: gameRepository,
		UserRepository: userRepository,
	}
}

func (gs *gameService) AddGameToUser(c context.Context, name,id string) error{
	// TODO TRANSACTIONS
	game,err:=gs.GameRepository.FindByName(c,name)
	if err!=nil{
		return err
	}
	user,err:=gs.UserRepository.FindById(c,id)
	if err!=nil{
		return err
	}
	game.NumberOfPlayers++
	user.Games=append(user.Games, game.Name)
	if err:=gs.GameRepository.Save(c,*game);err!=nil{
		return err
	}
	if err:=gs.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	return nil
}
func (gs *gameService) GetByName(c context.Context, name string) (*entity.Game,error){
	game,err:=gs.GameRepository.FindByName(c,name)
	if err!=nil{
		return nil,err
	}
	return game,nil
}
func (gs *gameService) GetByAmount(c context.Context, amount int) ([]entity.Game,error){
	games,err:=gs.GameRepository.FindByAmount(c,amount)
	if err!=nil{
		return nil,err
	}
	return games,nil
}
