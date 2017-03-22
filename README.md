# SurfNerd

A comprehensive suite of surf forecasting and wave analysis tools

### Documentation

https://godoc.org/github.com/mpiannucci/surfnerd

### What can it do

It can:

* Download data from NOAA WaveWatch 3 Model runs
* Download data frrom NOAA NAM Weather models
* Download buoy data from NOAA's vast buoy data base
* Find nearby buoys and model runs for given locations
* Find historical buoy data
* Solve a variety of wave equations to aide in wave height predictions and forecasts

### Are there examples of it being used? 

Yes! Check out [buoyfinder](https://buoyfinder.appspot.com) for a relatively straightforward usage of the buoy API.

### Why Go?

Most tools of this nature are used in MATLAB or Python. While I love both of those languages they each have flaws that turned me away. MATLAB is non-free (in both beer and speech) and a bit clunky. There is also no real great way to deploy MATLAB to the Google, AWS, etc cloud platforms. And while Python is well supported by all platforms, the free tier of Google App Engine (nor heroku) allow the use of c modules. So all of the scientific tooling I need to read from WaveWatch 3 (`hdf5`, `netcdf`) is not available. Which leads to speed. Because we are using iteration for grabbing a LOT of data instead of downloading the lean `GRIB` data ``Python is just too slow for the task.

Go is well supported by most if not all cloud platforms, is a drop dead simple language, is *free*, and has the easiest unit testing framework built in. Because of this I chose Go and it has worked decently well. In a perfect world Go will also allow for bindings to other languages such as c and Python so that will be something to work towards.