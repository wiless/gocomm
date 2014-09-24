package dsp

import (
	"github.com/mjibson/go-dsp/fft"
	"math"
	"math/cmplx"
	"runtime"
	"wiless/gocomm"
	"wiless/vlib"
)

func Db(linearValue float64) float64 {
	return 10.0 * math.Log10(linearValue)
}

func InvDb(dBValue float64) float64 {
	return math.Pow(10, dBValue/10.0)
}
func ExtFFT_C(samples vlib.VectorC, N int) vlib.VectorC {

	y := vlib.VectorC(fft.FFT(samples))
	y = y.Scale(math.Sqrt(float64(N)))
	return y
}
func ExtIFFT_C(samples vlib.VectorC, N int) vlib.VectorC {

	y := vlib.VectorC(fft.IFFT(samples))
	y = y.Scale(math.Sqrt(1.0 / float64(N)))
	return y

}

func FFT_C(samples vlib.VectorC, N int) vlib.VectorC {

	samples.Resize(N)
	fbins := vlib.NewVectorC(N)
	normalize := complex(math.Sqrt(1.0/float64(N)), 0)
	result := vlib.NewVectorC(N)
	for i := 0; i < N; i++ {
		for n := 0; n < N; n++ {
			scale := float64(i) * float64(n) / float64(N)
			binf := complex(0, -2.0*math.Pi*scale)
			fbins[n] = cmplx.Exp(binf)
		}
		// fbins = fbins.ScaleC(i)
		// fmt.Print("\ni=", i, fbins)
		result[i] = vlib.DotC(samples, fbins) * normalize
	}

	return result
}

func IFFT_C(samples vlib.VectorC, N int) vlib.VectorC {

	samples.Resize(N)
	fbins := vlib.NewVectorC(N)
	normalize := complex(math.Sqrt(1.0/float64(N)), 0)
	result := vlib.NewVectorC(N)
	for i := 0; i < N; i++ {
		for n := 0; n < N; n++ {
			scale := float64(i) * float64(n) / float64(N)
			binf := complex(0, 2.0*math.Pi*scale)
			fbins[n] = cmplx.Exp(binf)
		}
		// fbins = fbins.ScaleC(i)
		// fmt.Print("\ni=", i, fbins)
		result[i] = vlib.DotC(samples, fbins) * normalize
	}

	return result
}

func FFT(samples vlib.VectorF, N int) vlib.VectorC {
	var csamples vlib.VectorC
	csamples.SetVectorF(samples)
	return FFT_C(csamples, N)
}

func IFFT_F(samples vlib.VectorF, N int) vlib.VectorC {
	var csamples vlib.VectorC
	csamples.SetVectorF(samples)
	return IFFT_C(csamples, N)
}

func GoIFFT_C(samples vlib.VectorC, N int) vlib.VectorC {
	n := runtime.GOMAXPROCS(8)
	if N != samples.Size() {
		samples.Resize(N)
	}
	// fbins := vlib.NewVectorC(N)
	result := vlib.NewVectorC(N)
	NChannels := make([]gocomm.Complex128Channel, N)

	//bigChannel := make(gocomm.Complex128Channel, N)

	for i := 0; i < N; i++ {
		NChannels[i] = gocomm.NewComplex128Channel()
		go GoFFTPerK(NChannels[i], samples, i, N, true)
		//go GoFFTPerK(bigChannel, samples, i, N, false)
	}

	//for i := 0; i < N; i++ {
	//	result[i] = (<-bigChannel).Ch
	//}
	for i := 0; i < N; i++ {
		result[i] = (<-NChannels[i]).Ch
	}

	runtime.GOMAXPROCS(n)
	return result

}

func GoFFT_C(samples vlib.VectorC, N int) vlib.VectorC {
	n := runtime.GOMAXPROCS(8)
	if N != samples.Size() {
		samples.Resize(N)
	}
	// fbins := vlib.NewVectorC(N)
	result := vlib.NewVectorC(N)
	NChannels := make([]gocomm.Complex128Channel, N)

	//bigChannel := make(gocomm.Complex128Channel, N)

	for i := 0; i < N; i++ {
		NChannels[i] = gocomm.NewComplex128Channel()
		go GoFFTPerK(NChannels[i], samples, i, N, false)
		//go GoFFTPerK(bigChannel, samples, i, N, false)
	}

	//for i := 0; i < N; i++ {
	//	result[i] = (<-bigChannel).Ch
	//}
	for i := 0; i < N; i++ {
		result[i] = (<-NChannels[i]).Ch
	}

	runtime.GOMAXPROCS(n)
	return result

}

func GoFFTPerK(outputSymbol gocomm.Complex128Channel, inputsamples vlib.VectorC, k, N int, inverse bool) {

	kbyN := float64(k) / float64(N)
	normalize := complex(1.0/math.Sqrt(float64(N)), 0)
	if !inverse {
		kbyN = kbyN * -1.0
	}

	fbins := vlib.NewVectorC(N)

	for n := 0; n < N; n++ {
		scale := kbyN * float64(n)
		binf := complex(0, 2.0*math.Pi*scale)
		fbins[n] = cmplx.Exp(binf)
	}
	result := vlib.GoDotC(inputsamples, fbins, 4) * normalize
	// result := vlib.DotC(inputsamples, fbins) * normalize

	var data gocomm.SComplex128Obj
	data.Ch = result
	outputSymbol <- data
}
