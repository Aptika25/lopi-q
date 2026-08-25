export const POSKO_COORDS = {
  latitude: -5.5645,
  longitude: 120.1945,
  maxRadiusMeters: 2.0
}

export function calculateDistance(lat1: number, lon1: number, lat2 = POSKO_COORDS.latitude, lon2 = POSKO_COORDS.longitude): number {
  const R = 6371000 // Earth radius in meters
  const dLat = (lat2 - lat1) * (Math.PI / 180)
  const dLon = (lon2 - lon1) * (Math.PI / 180)
  
  const a = 
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(lat1 * (Math.PI / 180)) * Math.cos(lat2 * (Math.PI / 180)) * 
    Math.sin(dLon / 2) * Math.sin(dLon / 2)
    
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
  const distance = R * c
  return parseFloat(distance.toFixed(2))
}
