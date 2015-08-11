package neuralnet

import (
    "fmt"
    )

type Gride struct {
    InputsChan []chan float64
    InputLayer []Dispatcher
    HiddenLayers [][]Neurone
    ChanRegister [][][]chan float64
    OutputLayer []Summer
    OutputsChan []chan float64
}

func NewGride(layerLens []int) *Gride {
    
    G := Gride{}
    G.InputsChan = make([]chan float64, layerLens[0])
    G.OutputsChan = make([]chan float64, layerLens[len(layerLens)-1])
    
    G.ChanRegister = make([][][]chan float64, len(layerLens)-1)
    for i := 0;i<len(G.ChanRegister);i++ {
        G.ChanRegister[i] = make([][]chan float64, layerLens[i])
        for j := 0;j<len(G.ChanRegister[i]);j++ {
            G.ChanRegister[i][j] = make([]chan float64, layerLens[i+1])
            for k := 0;k<len(G.ChanRegister[i][j]);k++ {
                G.ChanRegister[i][j][k] = make(chan float64)
            }
        }
    }
    
    G.InputLayer = make([]Dispatcher, layerLens[0])
    for i, v := range G.InputLayer {
        
        v = Dispatcher{}
        v.LinkIn(G.InputsChan[i])
        for _, x := range G.ChanRegister[0][i] {
            v.LinkOutAdd(x)
        }
    }
    
    G.OutputLayer = make([]Summer, layerLens[len(layerLens)-1])
    
    for i, v := range G.OutputLayer {
        v = Summer{}
        v.LinkOut(G.OutputsChan[i])
        for j := 0;j<layerLens[len(layerLens)-2];j++ {
            v.LinkInAdd(G.ChanRegister[len(G.ChanRegister)-1][j][i])
        }
    }
    G.HiddenLayers = make([][]Neurone, len(layerLens)-2)
    for i, v := range G.HiddenLayers {
        v = make([]Neurone, layerLens[i+1])
        for j, k := range v {
            k = Neurone{}
            k.S = &Summer{}
            for n := 0;n<layerLens[i];n++ {
                k.S.LinkInAdd(G.ChanRegister[i][n][j])
            }
            k.chanST = make(chan float64)
            k.S.LinkOut(k.chanST)
            k.T = &Transitor{}
            k.T.LinkIn(k.chanST)
            k.chanTD = make(chan float64)
            k.T.LinkOut(k.chanTD)
            k.D = &Dispatcher{}
            k.D.LinkIn(k.chanTD)
            for _, x := range G.ChanRegister[i+1][j] {
                k.D.LinkOutAdd(x)
            }
            k.Start()
        }
    }
    return &G
}

func (s Gride) String() string {
    ret := fmt.Sprintf("Inputs:%d\n", len(s.InputLayer))
    for i := 0;i<len(s.ChanRegister);i++ {
        ret = fmt.Sprintf("%sLayer %d: \n", ret, i)
        for j := 0;j<len(s.ChanRegister[i]);j++ {
            ret = fmt.Sprintf("%s\tNode %d: \n",ret,j)
            for k := 0;k<len(s.ChanRegister[i][j]);k++ {
                ret = fmt.Sprintf("%s\t\tchannel %d\n",ret,k)
            }
        }
    }
    ret = fmt.Sprintf("%sOutputs:%d\n", ret, len(s.OutputLayer))
    return ret
}

func (s Gride)Push(inputTab []float64) {
    for i, v := range s.InputsChan {
        v <- inputTab[i]
    }
}
func (s Gride)Get() (ret []float64) {
    for i, v := range s.OutputsChan {
        ret[i] = <-v
    }
    return
}