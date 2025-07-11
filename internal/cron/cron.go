package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
	"github.com/t3mp14r3/shiny-umbrella/internal/repository"
	"go.uber.org/zap"
)

type Cron struct {
    scheduler   *gocron.Scheduler
    repo        *repository.Repository
    logger      *zap.Logger
    jobs        []*gocron.Job
    scores      []*gocron.Job
}

func New(repo *repository.Repository, logger *zap.Logger) *Cron {
    s, err := gocron.NewScheduler()
	
    if err != nil {
        log.Fatalf("Failed to initialize new cron scheduler: %v\n", err)
        return nil
	}

    return &Cron{
        scheduler: &s,
        repo: repo,
        logger: logger,
    }
}

func (c *Cron) Start() {
    (*c.scheduler).Start()
}

func (c *Cron) Close() error {
    return (*c.scheduler).Shutdown()
}

func (c *Cron) Load() error {
    tournaments, err := c.repo.GetTournaments(context.Background())

    if err != nil {
        return err
    }

    for i, t := range tournaments {
        endsAt := t.StartsAt.Add(time.Duration(t.Duration * int64(time.Second)))
        
        if time.Now().Before(endsAt) {
            err = c.Score(&tournaments[i])
        
            if err != nil {
                return err
            }
        }
    }
    log.Println("Loaded regular")

    records, err := c.repo.GetAutomatics(context.Background())

    if err != nil {
        return err
    }

    log.Println(records)

    for _, record := range records {
        j, err := (*c.scheduler).NewJob(
            gocron.DurationJob(
                time.Duration(record.Repeat*int64(time.Second)),
            ),
            gocron.NewTask(
                func(input domain.Tournament, rewards []uint8, logger *zap.Logger) {
                    log.Println("lets goooo")
                    input.StartsAt = time.Now()
                    result, err := c.repo.CreateTournament(context.Background(), input)
        
                    if err != nil {
                        return
                    }

                    var rew []domain.Reward
                    err = json.Unmarshal([]byte(rewards), &rew)
                    
                    if err != nil {
                        c.logger.Error("Failed to unmarshal rewards data", zap.Error(err))
                        return
                    }

                    err = c.repo.CreateRewards(context.Background(), result.ID, rew)
                    
                    if err != nil {
                        return
                    }
                },
                domain.Tournament{
                    Price: record.Price,
                    MinUsers: record.MinUsers,
                    MaxUsers: record.MaxUsers,
                    Bets: record.Bets,
                    Duration: record.Duration,
                },
                record.Rewards,
                c.logger,
            ),
            gocron.WithName(fmt.Sprintf("%d", record.ID)),
            gocron.WithStartAt(gocron.WithStartImmediately()),
        )

        if err != nil {
            c.logger.Error("Failed to start a new cron job", zap.Error(err))
            continue
        }

        c.jobs = append(c.jobs, &j)
    }
    log.Println("Loaded automatics")



    log.Println(c.jobs)

    return nil
}

func (c *Cron) Score(input *domain.Tournament) error {
    endsAt := input.StartsAt.Add(time.Duration(input.Duration * int64(time.Second)))

    log.Println("scheduling new score!")

    go func(id int64) {
        time.Sleep(time.Until(endsAt))
        c.repo.Calculate(context.Background(), id)
    }(input.ID)

    return nil
}

func (c *Cron) Update(id int64, channel string) error {
    record, err := c.repo.GetAutomatic(context.Background(), id)

    if err != nil {
        c.logger.Error("Failed to get automatic record", zap.Error(err))
        return err
    }

    if channel == "inserts" {
        log.Println("new insert!")
        log.Println(fmt.Sprintf("%d", record.ID))

        j, err := (*c.scheduler).NewJob(
            gocron.DurationJob(
                time.Duration(record.Repeat*int64(time.Second)),
            ),
            gocron.NewTask(
                func(input domain.Tournament, rewards []uint8, logger *zap.Logger) {
                    input.StartsAt = time.Now()
                    result, err := c.repo.CreateTournament(context.Background(), input)
        
                    if err != nil {
                        return
                    }

                    var rew []domain.Reward
                    err = json.Unmarshal([]byte(rewards), &rew)
                    
                    if err != nil {
                        c.logger.Error("Failed to unmarshal rewards data", zap.Error(err))
                        return
                    }

                    err = c.repo.CreateRewards(context.Background(), result.ID, rew)
                    
                    if err != nil {
                        return
                    }
                },
                domain.Tournament{
                    Price: record.Price,
                    MinUsers: record.MinUsers,
                    MaxUsers: record.MaxUsers,
                    Bets: record.Bets,
                    Duration: record.Duration,
                },
                record.Rewards,
                c.logger,
            ),
            gocron.WithName(fmt.Sprintf("%d", record.ID)),
            gocron.WithStartAt(gocron.WithStartImmediately()),
        )
        
        if err != nil {
            c.logger.Error("Failed to start a new cron job", zap.Error(err))
            return err
        }

        c.jobs = append(c.jobs, &j)
    } else if channel == "updates" {
        var jobID uuid.UUID
        for i, job := range c.jobs {
            if (*job).Name() == fmt.Sprintf("%d", id) {
                jobID = (*job).ID()
                c.jobs = append(c.jobs[:i], c.jobs[i+1:]...)
                break
            }
        }
        (*c.scheduler).RemoveJob(jobID)
        
        j, err := (*c.scheduler).NewJob(
            gocron.DurationJob(
                time.Duration(record.Repeat*int64(time.Second)),
            ),
            gocron.NewTask(
                func(input domain.Tournament, rewards []uint8, logger *zap.Logger) {
                    input.StartsAt = time.Now()
                    result, err := c.repo.CreateTournament(context.Background(), input)
        
                    if err != nil {
                        return
                    }

                    var rew []domain.Reward
                    err = json.Unmarshal([]byte(rewards), &rew)
                    
                    if err != nil {
                        c.logger.Error("Failed to unmarshal rewards data", zap.Error(err))
                        return
                    }

                    err = c.repo.CreateRewards(context.Background(), result.ID, rew)
                    
                    if err != nil {
                        return
                    }
                },
                domain.Tournament{
                    Price: record.Price,
                    MinUsers: record.MinUsers,
                    MaxUsers: record.MaxUsers,
                    Bets: record.Bets,
                    Duration: record.Duration,
                },
                record.Rewards,
                c.logger,
            ),
            gocron.WithName(fmt.Sprintf("%d", record.ID)),
            gocron.WithStartAt(gocron.WithStartImmediately()),
        )
        
        if err != nil {
            c.logger.Error("Failed to start an updated cron job", zap.Error(err))
            return err
        }

        c.jobs = append(c.jobs, &j)
    } else if channel == "deletes" {
        var jobID uuid.UUID
        for i, job := range c.jobs {
            if (*job).Name() == fmt.Sprintf("%d", id) {
                jobID = (*job).ID()
                c.jobs = append(c.jobs[:i], c.jobs[i+1:]...)
                break
            }
        }
        (*c.scheduler).RemoveJob(jobID)
    }
    
    return nil
}

func (c *Cron) Regular(id int64) error {
    record, err := c.repo.GetTournament(context.Background(), id)

    if err != nil {
        return err
    }

    return c.Score(record)
}
