package main

import (
	"fmt"
	"time"

	"github.com/airbusgeo/godal"
)

func main() {
	t := time.Now()
	t1 := t
	// redfile := "/data/Projects/Go/test/S2A_MSIL2A_20190531T081611_N0212_R121_T37TEM_20190531T112906.SAFE/GRANULE/L2A_T37TEM_A020567_20190531T081627/IMG_DATA/R10m/T37TEM_20190531T081611_B08_10m.jp2"
	// nirfile := "/data/Projects/Go/test/S2A_MSIL2A_20190531T081611_N0212_R121_T37TEM_20190531T112906.SAFE/GRANULE/L2A_T37TEM_A020567_20190531T081627/IMG_DATA/R10m/T37TEM_20190531T081611_B04_10m.jp2"
	redfile := "T37TEM_20190531T081611_B08_10m.jp2"
	nirfile := "T37TEM_20190531T081611_B04_10m.jp2"

	godal.RegisterAll()
	ds_red, err := godal.Open(redfile)
	if err != nil {
		panic(err)
	}
	defer ds_red.Close()
	fmt.Print("Открываем красный - ")
	red_structure := ds_red.Structure()
	// fmt.Printf("Size is %dx%dx%d\n", red_structure.SizeX, red_structure.SizeY, red_structure.NBands)
	red_band := ds_red.Bands()[0]
	red_pafScanline := make([]float32, red_structure.SizeX*red_structure.SizeY)
	err = red_band.Read(0, 0, red_pafScanline, red_structure.SizeX, red_structure.SizeY)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(t1))
	t1 = time.Now()

	fmt.Print("Открываем инфракрасный - ")
	ds_nir, err := godal.Open(nirfile)
	if err != nil {
		panic(err)
	}
	defer ds_nir.Close()
	nir_structure := ds_nir.Structure()
	// fmt.Printf("Size is %dx%dx%d\n", nir_structure.SizeX, nir_structure.SizeY, nir_structure.NBands)
	nir_band := ds_nir.Bands()[0]

	nir_pafScanline := make([]float32, nir_structure.SizeX*nir_structure.SizeY)
	err = nir_band.Read(0, 0, nir_pafScanline, nir_structure.SizeX, nir_structure.SizeY)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(t1))
	t1 = time.Now()

	fmt.Print("Calculate NDVI - ")
	ndvi := make([]float32, nir_structure.SizeX*nir_structure.SizeY)
	for i := range red_pafScanline {
		ndvi[i] = 100 * (nir_pafScanline[i] - red_pafScanline[i]) / (nir_pafScanline[i] + red_pafScanline[i])
	}
	fmt.Println(time.Since(t1))
	t1 = time.Now()
	fmt.Print("Сохраняем GEOTiff - ")

	new, err := godal.Create(godal.GTiff, "test1.tiff", 1, godal.Int32, red_structure.SizeX, red_structure.SizeY, godal.CreationOption("COMPRESS=DEFLATE", "TILED=YES"))
	gt, _ := ds_red.GeoTransform()
	new.SetGeoTransform(gt)
	new.SetProjection(ds_red.Projection())
	new.Write(0, 0, ndvi, red_structure.SizeX, red_structure.SizeY)
	new.Close()
	fmt.Println(time.Since(t1))

	fmt.Println("Обсчет NDVI занял ", time.Since(t))

}
