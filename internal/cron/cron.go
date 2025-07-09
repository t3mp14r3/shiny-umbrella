package cron

import (
	"context"
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
}

func New(repo *repository.Repository, logger *zap.Logger) (*Cron, error) {
    s, err := gocron.NewScheduler()
	
    if err != nil {
        log.Printf("Failed to initialize new cron scheduler: %v\n", err)
        return nil, err
	}

    return &Cron{
        scheduler: &s,
        repo: repo,
        logger: logger,
    }, nil
}

func (c *Cron) Start() {
    (*c.scheduler).Start()
}

func (c *Cron) Shutdown() error {
    return (*c.scheduler).Shutdown()
}

func (c *Cron) Load() error {
    records, err := c.repo.GetAutomatics(context.Background())

    if err != nil {
        log.Printf("Failed to load automatic records: %v\n", err)
        return err
    }

    for _, record := range records {
        j, err := (*c.scheduler).NewJob(
            gocron.DurationJob(
                time.Duration(record.Repeat*int64(time.Second)),
            ),
            gocron.NewTask(
                func(input domain.Tournament, logger *zap.Logger) {
                    c.repo.CreateTournament(context.Background(), input)
                },
                domain.Tournament{
                    Price: record.Price,
                    MinUsers: record.MinUsers,
                    MaxUsers: record.MaxUsers,
                    Bets: record.Bets,
                    StartsAt: record.StartsAt,
                    Duration: record.Duration,
                },
                c.logger,
            ),
            gocron.WithName(fmt.Sprintf("%d", record.ID)),
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

func (c *Cron) Update(id int64, channel string) error {
    record, err := c.repo.GetAutomatic(context.Background(), id)

    if err != nil {
        c.logger.Error("Failed to get automatic record", zap.Error(err))
        return err
    }

    if channel == "inserts" {
        _, err := c.repo.CreateTournament(context.Background(), domain.Tournament{
            Price: record.Price,
            MinUsers: record.MinUsers,
            MaxUsers: record.MaxUsers,
            Bets: record.Bets,
            StartsAt: record.StartsAt,
            Duration: record.Duration,
        })

        if err != nil {
            return err
        }
        
        j, err := (*c.scheduler).NewJob(
            gocron.DurationJob(
                time.Duration(record.Repeat*int64(time.Second)),
            ),
            gocron.NewTask(
                func(input domain.Tournament, logger *zap.Logger) {
                    c.repo.CreateTournament(context.Background(), input)
                },
                domain.Tournament{
                    Price: record.Price,
                    MinUsers: record.MinUsers,
                    MaxUsers: record.MaxUsers,
                    Bets: record.Bets,
                    StartsAt: record.StartsAt,
                    Duration: record.Duration,
                },
                c.logger,
            ),
            gocron.WithName(fmt.Sprintf("%d", record.ID)),
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
                func(input domain.Tournament, logger *zap.Logger) {
                    c.repo.CreateTournament(context.Background(), input)
                },
                domain.Tournament{
                    Price: record.Price,
                    MinUsers: record.MinUsers,
                    MaxUsers: record.MaxUsers,
                    Bets: record.Bets,
                    StartsAt: record.StartsAt,
                    Duration: record.Duration,
                },
                c.logger,
            ),
            gocron.WithName(fmt.Sprintf("%d", record.ID)),
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
