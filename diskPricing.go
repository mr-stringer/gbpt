package main

import "log"

type Sploc struct {
	sp  string
	loc string
}

func (c Config) ReduceDisks() []DiskPrice {
	/* The purpose here is to find all of the different disks in different */
	/* locations we need and reduce them */

	/* collect unique combinations of storage profiles and locations */
	sp1 := []Sploc{}

	for _, sp := range c.StorageProfiles {
		for _, ap := range c.Applications {
			for _, env := range ap.Environments {
				for _, vm := range env.VMs {
					if sp.Name == vm.StorageProfile {
						var match = false
						for _, v := range sp1 {
							if v.sp == sp.Name && env.Location == v.loc {
								match = true
							}
						}
						if !match {
							sp1 = append(sp1, Sploc{sp.Name, env.Location})
						}
					}
				}
			}
		}
	}
	log.Printf("Hopefully print unique combos of StorageProfiles and Locations")
	for _, v := range sp1 {
		log.Printf("StorageProfile:%s, Location:%s", v.sp, v.loc)
	}
	/* Now we have the unique storage profile and location combos */
	/* creating unique disk/location combos should be easier */
	rd1 := []DiskPrice{}

	for _, sp := range c.StorageProfiles {
		for _, sploc := range sp1 {
			if sploc.sp != sp.Name {
				/* if the storage profiles don't match do a skip */
				continue
			}
			for _, disk := range sp.Disks {
				dp := DiskPrice{}
				dp.Location = sploc.loc
				if disk.Type == "pssd" {
					dp.DiskType = pssd
					dp.Pssd = GetPssdFromSize(disk.Size)
				} else {
					dp.DiskType = pssdv2
				}
				log.Printf("DiskType:%d pssdType:%s, Location:%s", dp.DiskType, dp.Pssd, dp.Location)

				match := false
				/* Check for a match in the returning slice */
				for _, rd := range rd1 {
					if rd.Location == dp.Location && rd.DiskType == dp.DiskType && rd.Pssd == dp.Pssd {
						match = true
					}
				}
				if !match {
					log.Print("Adding combo")
					rd1 = append(rd1, dp)
				}

			}
		}
	}
	log.Printf("Hopefully print a list of unique disks - found %d configs", len(rd1))
	for _, rd := range rd1 {
		log.Printf("DiskType:%d pssdType:%s, Location:%s", rd.DiskType, rd.Pssd, rd.Location)
	}

	return []DiskPrice{}
}
