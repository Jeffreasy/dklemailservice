package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"time"
)

type NewsletterService struct {
	fetcher   *NewsletterFetcher
	processor *NewsletterProcessor
	formatter *NewsletterFormatter
	sender    *NewsletterSender
	ticker    *time.Ticker
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewNewsletterService(fetcher *NewsletterFetcher, processor *NewsletterProcessor,
	formatter *NewsletterFormatter, sender *NewsletterSender) *NewsletterService {
	ticker := time.NewTicker(24 * time.Hour)
	ctx, cancel := context.WithCancel(context.Background())
	return &NewsletterService{
		fetcher:   fetcher,
		processor: processor,
		formatter: formatter,
		sender:    sender,
		ticker:    ticker,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (ns *NewsletterService) Start() {
	if ns == nil {
		return
	}
	go func() {
		for {
			select {
			case <-ns.ctx.Done():
				ns.ticker.Stop()
				return
			case <-ns.ticker.C:
				if err := ns.RunPipeline(); err != nil {
					logger.Error("Nieuwsbrief pipeline error", "error", err)
					if ns.sender != nil && ns.sender.notifSvc != nil {
						ns.sender.notifSvc.CreateNotification(ns.ctx, models.NotificationTypeSystem, models.NotificationPriorityMedium, "Nieuwsbrief Fout", err.Error())
					}
				}
			}
		}
	}()
	logger.Info("Newsletter service gestart")
}

func (ns *NewsletterService) RunPipeline() error {
	raw, err := ns.fetcher.Fetch(ns.ctx)
	if err != nil {
		return err
	}
	processed := ns.processor.Process(raw)
	content, err := ns.formatter.Format(&processed, "DKL Wekelijkse Nieuwsbrief")
	if err != nil {
		return err
	}
	return ns.sender.Send(ns.ctx, content, "DKL Wekelijkse Nieuwsbrief")
}

func (ns *NewsletterService) Stop() {
	if ns != nil {
		ns.cancel()
	}
}
