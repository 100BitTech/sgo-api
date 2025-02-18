package base

import (
	"sync"

	"github.com/samber/oops"
)

// 简单快速的并发执行，需先指定并发数量、消费者、一组生产元素
func Go[I any, O any](
	concurrent int,
	consume func(I) (O, error),
	items ...I,
) ([]O, error) {
	oCh := make(chan O, len(items)+1)
	errCh := make(chan error, len(items)+1)

	ch := make(chan I, concurrent*2)
	var wg sync.WaitGroup
	for i := 0; i < concurrent; i++ {
		wg.Add(1)

		go func() {
			Try(func() {
				for i := range ch {
					o, err := consume(i)

					if err != nil {
						errCh <- oops.Wrap(err)
					} else {
						oCh <- o
					}
				}
			}).Catch(func(err error) {
				errCh <- oops.Wrapf(err, "恐慌")
			}).Finally(func() {
				wg.Done()
			}).Do()
		}()
	}

	for _, i := range items {
		ch <- i
	}
	close(ch)
	wg.Wait()
	close(oCh)
	close(errCh)

	if len(errCh) > 0 {
		errs := []error{}
		for e := range errCh {
			errs = append(errs, e)
		}
		return nil, oops.Wrap(NewMultiError(errs...))
	} else {
		os := []O{}
		for o := range oCh {
			os = append(os, o)
		}
		return os, nil
	}
}
