package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"github.com/montanaflynn/stats"
)

// linear regression function
func linearRegression(x, y []float64) (slope, intercept float64, err error) {
	if len(x) != len(y) {
		return 0, 0, fmt.Errorf("x and y must have the same length")
	}
	n := float64(len(x))

	sumX, _ := stats.Sum(x)
	sumY, _ := stats.Sum(y)
	sumXY := 0.0
	sumX2 := 0.0

	for i := 0; i < len(x); i++ {
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	numerator := n*sumXY - sumX*sumY
	denominator := n*sumX2 - sumX*sumX

	slope = numerator / denominator
	meanX, _ := stats.Mean(x)
	meanY, _ := stats.Mean(y)
	intercept = meanY - slope*meanX

	return slope, intercept, nil
}

func regressionMetrics(x, y []float64, slope, intercept float64) (r2, rse, f float64) {
	n := float64(len(x))

	yMean, _ := stats.Mean(y)

	var rss, tss float64
	for i := 0; i < len(x); i++ {
		yHat := slope*x[i] + intercept
		residual := y[i] - yHat
		rss += residual * residual

		tss += (y[i] - yMean) * (y[i] - yMean)
	}

	r2 = 1 - rss/tss
	rse = math.Sqrt(rss / (n - 2))
	f = ((tss - rss) / 1) / (rss / (n - 2))

	return
}

func makePlot(title string, x, y []float64, slope, intercept float64) *plot.Plot {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.X.Min = 2
	p.X.Max = 20
	p.Y.Min = 2
	p.Y.Max = 14

	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	scatter, _ := plotter.NewScatter(pts)
	scatter.GlyphStyle.Radius = vg.Points(3)
	p.Add(scatter)

	lineData := plotter.XYs{
		{X: 2, Y: slope*2 + intercept},
		{X: 20, Y: slope*20 + intercept},
	}
	line, _ := plotter.NewLine(lineData)
	p.Add(line)

	return p
}

func main() {
	// define Anscombe Quartet data
	x1 := []float64{10, 8, 13, 9, 11, 14, 6, 4, 12, 7, 5}
	x2 := []float64{10, 8, 13, 9, 11, 14, 6, 4, 12, 7, 5}
	x3 := []float64{10, 8, 13, 9, 11, 14, 6, 4, 12, 7, 5}
	x4 := []float64{8, 8, 8, 8, 8, 8, 8, 19, 8, 8, 8}
	y1 := []float64{8.04, 6.95, 7.58, 8.81, 8.33, 9.96, 7.24, 4.26, 10.84, 4.82, 5.68}
	y2 := []float64{9.14, 8.14, 8.74, 8.77, 9.26, 8.1, 6.13, 3.1, 9.13, 7.26, 4.74}
	y3 := []float64{7.46, 6.77, 12.74, 7.11, 7.81, 8.84, 6.08, 5.39, 8.15, 6.42, 5.73}
	y4 := []float64{6.58, 5.76, 7.71, 8.84, 8.47, 7.04, 5.25, 12.5, 5.56, 7.91, 6.89}

	// list of datasets
	datasets := []struct {
		name string
		x    []float64
		y    []float64
	}{
		{"Set I", x1, y1},
		{"Set II", x2, y2},
		{"Set III", x3, y3},
		{"Set IV", x4, y4},
	}

	var totalDuration time.Duration
	var totalMemoryUsed uint64

	// loop through datasets
	for i, data := range datasets {
		var mStart, mEnd runtime.MemStats

		runtime.GC() // clean before measuring

		// capture memory before execution
		runtime.ReadMemStats(&mStart)
		start := time.Now()

		slope, intercept, err := linearRegression(data.x, data.y)
		if err != nil {
			fmt.Printf("%s: regression error: %v\n", data.name, err)
			continue
		}

		duration := time.Since(start)

		runtime.GC() //clean after measuring
		// capture memory after execution
		runtime.ReadMemStats(&mEnd)

		memoryUsed := mEnd.Alloc - mStart.Alloc

		totalDuration += duration
		totalMemoryUsed += memoryUsed

		finalAlloc := mEnd.Alloc // just record the ending heap allocation

		// using "memUsed := mEnd.Alloc - mStart.Alloc":
		// if mEnd.Alloc < mStart.Alloc, it becomes a negative number
		// and the "memUsed" will go around in a circle and becomes an oversized number

		r2, rse, fstat := regressionMetrics(data.x, data.y, slope, intercept)

		fmt.Printf("%s:\n", data.name)
		fmt.Printf("  slope = %.5f\n", slope)
		fmt.Printf("  intercept = %.5f\n", intercept)

		fmt.Printf("  R-squared = %.4f\n", r2)
		fmt.Printf("  Residual Std Error = %.4f\n", rse)
		fmt.Printf("  F-statistic = %.4f\n", fstat)

		fmt.Printf("  execution time = %v\n", duration)
		fmt.Printf("  Memory used = %d bytes (%.2f KB)\n\n", finalAlloc, float64(finalAlloc)/1024.0)

		p := makePlot(data.name, data.x, data.y, slope, intercept)

		filename := fmt.Sprintf("anscombe_go_set_%d.png", i+1)
		if err := p.Save(5*vg.Inch, 5*vg.Inch, filename); err != nil {
			log.Fatalf("Failed to save %s: %v", filename, err)
		}
	}
	fmt.Println("--------------------------------------------------")
	fmt.Println("Summary for All Sets:")
	fmt.Printf("  Total Execution Time = %v\n", totalDuration)
	fmt.Printf("  Total Memory Used = %d bytes (%.2f KB)\n", totalMemoryUsed, float64(totalMemoryUsed)/1024.0)
	fmt.Println("--------------------------------------------------")

}

/*
Set I:
slope = 0.50009
intercept = 3.00009
execution time = 10.709µs
Memory used = 120328 bytes (117.51 KB)

Set II:
slope = 0.50000
intercept = 3.00091
execution time = 250ns
Memory used = 125688 bytes (122.74 KB)

Set III:
slope = 0.49973
intercept = 3.00245
execution time = 209ns
Memory used = 125704 bytes (122.76 KB)

Set IV:
slope = 0.49991
intercept = 3.00173
execution time = 167ns
Memory used = 125704 bytes (122.76 KB)
*/
// I noticed that it used >100kb memories, why?
// So even if my dataset is only a few dozen float64s and it would theoretically only need about 1KB of storage.
// a Go program will cost at least 100kb to start with due to the runtime + packages + memory pooling mechanism.
// ---> the 120kb is not memory for data, but for the environment that the program to run
