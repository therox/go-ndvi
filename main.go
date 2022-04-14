package main

import (
	"fmt"
	"time"

	"github.com/airbusgeo/godal"
)

func main() {
	t := time.Now()
	t1 := t

	redfile := "sentinel_example_data/T37UDU_20171111T085159_B08_10m.jp2"
	nirfile := "sentinel_example_data/T37UDU_20171111T085159_B04_10m.jp2"

	godal.RegisterAll()
	ds_red, err := godal.Open(redfile)
	if err != nil {
		panic(err)
	}
	defer ds_red.Close()
	fmt.Print("Open red channel - ")
	red_structure := ds_red.Structure()
	red_band := ds_red.Bands()[0]
	red_pafScanline := make([]float32, red_structure.SizeX*red_structure.SizeY)
	err = red_band.Read(0, 0, red_pafScanline, red_structure.SizeX, red_structure.SizeY)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(t1))
	t1 = time.Now()

	fmt.Print("Open infrared channel - ")
	ds_nir, err := godal.Open(nirfile)
	if err != nil {
		panic(err)
	}
	defer ds_nir.Close()

	nir_structure := ds_nir.Structure()
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
	fmt.Print("Saving GEOTiff - ")

	new, err := godal.Create(godal.GTiff, "test1.tiff", 1, godal.Int32, red_structure.SizeX, red_structure.SizeY, godal.CreationOption("COMPRESS=DEFLATE", "TILED=YES"))
	gt, _ := ds_red.GeoTransform()
	new.SetGeoTransform(gt)
	new.SetProjection(ds_red.Projection())
	new.Write(0, 0, ndvi, red_structure.SizeX, red_structure.SizeY)
	new.Close()
	fmt.Println(time.Since(t1))

	fmt.Println("NDVI calculated in ", time.Since(t))

}
