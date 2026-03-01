package geoIP

func (r *GeoIP) IP2Location(ipAddr string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results, err := r.db.Get_all(ipAddr)
	if err == nil && results.City != "" {
		return results.City, nil
	}

	if err == nil && results.Region != "" {
		return results.Region, nil
	}

	if err == nil && results.Country_short != "" {
		return results.Country_short, nil
	}

	return "", err
}
