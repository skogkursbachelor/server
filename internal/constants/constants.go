package constants

// DefaultPath for the server
const DefaultPath = "/"

// Version of the API
const Version = "v1"

const APIPath = DefaultPath + "api/" + Version + "/"
const ForestryRoadsPath = APIPath + "forestryroads"

const ProxyPath = DefaultPath + "proxy/"

// API endpoints

const MapTilerTransformAPI = "https://api.maptiler.com/coordinates/transform/"
const NVEFrostDepthAPI = "https://gts.nve.no/api/MultiPointTimeSeries/ByMapCoordinateCsv"
