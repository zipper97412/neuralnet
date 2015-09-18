package neuralnet

import (
    "fmt"
    )
    
type Transitor struct {
    InputChan <-chan float64
    OutputChan chan<- float64
}
func (t *Transitor)LinkOut(ch chan<- float64) {
    t.OutputChan = ch
}
func (t *Transitor)LinkIn(ch <-chan float64) {
    t.InputChan = ch
}


func (t *Transitor)Start() {
    go func() {
        for {
            input, more := <-t.InputChan
            t.OutputChan <- t.transit( input )
            if !more {
                close(t.OutputChan)
                return
            }
        
        }
    }()
}

func (t *Transitor)transit( input float64) float64 {
    
    fmt.Printf("in transit %f",input)
    return 1-input
}




type Summer struct {
    InputsChan []<-chan float64
    InputsTemp []float64
    Weights []float64
    OutputChan chan<- float64
    updated chan bool
}

func (s *Summer)LinkInAdd(ch <-chan float64) {
    s.InputsChan = append(s.InputsChan, ch)
    s.Weights = append(s.Weights, 1)
    index := len(s.InputsChan)-1
    go func(i int) {
        input, more := <-s.InputsChan[i]
        s.InputsTemp[i] = input
        s.updated <- true
        if !more {
            if len(s.InputsChan) == 1 {
                close(s.updated)
                return
            }
            s.InputsChan = append(s.InputsChan[:i], s.InputsChan[i+1:]...)
            s.Weights = append(s.Weights[:i], s.Weights[i+1:]...)
            return
        }
    }(index)
}

func (s *Summer)LinkOut(ch chan<- float64) {
    s.OutputChan = ch
}
    
    

func (s *Summer)Start() {
    s.updated = make(chan bool, len(s.InputsChan))
    go func() {
        count := 0
        for {
            _, more := <-s.updated
            count++
            if count >= len(s.InputsChan) {
                s.OutputChan <- s.summ()
            }
            if !more {
                close(s.OutputChan)
                return
            }
        }
    }()
}

func (s Summer)summ() float64 {
    ret := 0.0
    for i := 0;i<len(s.InputsTemp);i++ {
        ret = ret + s.Weights[i]*s.InputsTemp[i]
    }
    return ret/float64(len(s.InputsTemp))
}





type Dispatcher struct {
    InputChan <-chan float64
    OutputsChan []chan<- float64
}

func (s *Dispatcher)LinkIn(ch <-chan float64) {
    s.InputChan = ch
}


func (s *Dispatcher)LinkOutAdd(ch chan<- float64) {
    s.OutputsChan = append(s.OutputsChan, ch)
}


func (s *Dispatcher)Start() {
    
    go func() {
        fmt.Println("dispatcher")
        input, more := <- s.InputChan
        fmt.Println("réception de",input, "dans Dispatcher")
        for _, v := range s.OutputsChan {
            v <- input
            fmt.Println("envoi de",input, "par Dispatcher")
        }
        if !more {
            for _, v := range s.OutputsChan {
                close(v)
            }
        }
    }()
    
}
    
    
    
type Neurone struct {
    S *Summer
    chanST chan float64
    T *Transitor
    chanTD chan float64
    D *Dispatcher
}
func (s *Neurone)Start() {
    s.S.Start()
    s.T.Start()
    s.D.Start()
}
func NewNeurone() Neurone {
    ret := Neurone{}
    ret.S = &Summer{}
    ret.T = &Transitor{}
    ret.D = &Dispatcher{}
    ret.chanTD = make(chan float64)
    ret.chanST = make(chan float64)
    ret.S.LinkOut(ret.chanST)
    ret.T.LinkIn(ret.chanST)
    ret.T.LinkOut(ret.chanTD)
    ret.D.LinkIn(ret.chanTD)
    return ret
}
