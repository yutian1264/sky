/*
@Time : 2018/10/25 14:17 
@Author : sky
@Description:
@File : re
@Software: GoLand
图1 波浪曲面

这个曲面的函数是：

z=f(x,y)=10×sinx2+y2−−−−−−√x2+y2−−−−−−√

不妨令
r=x2+y2−−−−−−√
，则

z=sin(r)r

2. 原理和实现

下面代码路径是：src/gopl/basictypes/wave.

下面的实现里，为了能让图形更好的适应画布，使用了一些比例缩放。

    先假设有 100×100

个小方格，接下来假设这 100 个小方格的 x 坐标变化范围是 [−15,15]， y 坐标范围和 x
一样。
通过 f
函数计算出曲面的高度 z
.
将小方格的第一个顶点 (x,y,z)
投影到二维坐标系中，得到坐标 (xx,yy)
。（使用 project 函数）
将 (xx,yy)
映射到画布中合适的位置。
接下来用同样的方法，将小方格的另外三个顶点映射到画布中合适的位置。
计算出小方格4个顶点在画布中的合适位置后，使用 svg 的 polygon 指令绘制出来。




*/
package main

import (
	"math"
	"io"

	"net/http"
	"fmt"
)

const (
	angle   = math.Pi / 6
	width   = 1000
	height  = 1000
	xyrange = 30.0 // x, y  变化范围，2*pi 一个周期，大约 2*2.4 个周期
	xscale  = width / 2 / xyrange
	yscale  = height / 2 / xyrange
	cells   = 100
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func draw(w io.Writer) {
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' style='stroke: red; fill: white; stroke-width: 0.7' width='1000' height='1000'>\n")

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			xx1, yy1 := corner(i, j)
			xx2, yy2 := corner(i, j+1)
			xx3, yy3 := corner(i+1, j+1)
			xx4, yy4 := corner(i+1, j)
			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g' />\n", xx1, yy1, xx2, yy2, xx3, yy3, xx4, yy4)
		}
	}
	fmt.Fprintf(w, "</svg>\n")
}

// 计算第 (i, j) 个小方块的角点坐标
func corner(i, j int) (float64, float64) {
	// 将 i, j 映射到 xyrange 这个范围，然后再计算坐标
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	xx, yy := project(x, y, f(x, y))
	return width/2 + xx*xscale, height/2 + yy*yscale
}
func project(x, y, z float64) (float64, float64) {
	xx := x*cos30 - y*cos30
	yy := x*sin30 + y*sin30 - z
	return xx, yy
}
func f(x float64, y float64) float64 {
	r := math.Hypot(x, y)
	z := math.Sin(r) / r
	return 10 * z
}
func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	draw(w)
}
func main() {
	http.HandleFunc("/", handle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
