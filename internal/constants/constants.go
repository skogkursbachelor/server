package constants

// DefaultPath for the server
const DefaultPath = "/"

// Version of the API
const Version = "v1"

// Internal Endpoints

const APIPath = DefaultPath + "api/" + Version + "/"
const ForestryRoadsPath = APIPath + "forestryroads"

const ProxyPath = DefaultPath + "proxy/"

// External Endpoints

const NVEFrostDepthAPI = "https://gts.nve.no/api/MultiPointTimeSeries/ByMapCoordinateCsv"
const NVEAreaTimeSeriesAPI = "https://gts.nve.no/api/GridTimeSeries/{x}/{y}/{startdate}/{enddate}/{theme}.{format}"
const ForestryRoadsWFS = "https://wms.geonorge.no/skwms1/wms.traktorveg_skogsbilveger"
const OpenMeteoEnsembleAPI = "https://ensemble-api.open-meteo.com/v1/ensemble?latitude={latitude}&longitude={longitude}&hourly=soil_moisture_10_to_40cm&models=gfs05&start_date={start_date}&end_date={end_date}"
const OpenMeteoDeepSoilTempURL = "https://api.open-meteo.com/v1/forecast?latitude={latitude}&longitude={longitude}&hourly=soil_temperature_54cm&start_date={start_date}&end_date={end_date}&models=icon_seamless"
const OpenMeteoHistoricalDeepSoilTempURL = "https://historical-forecast-api.open-meteo.com/v1/forecast?latitude={latitude}&longitude={longitude}&start_date={start_date}&end_date={end_date}&hourly=soil_temperature_54cm&models=icon_seamless"
