package fs

import "time"

func GetEconetData(date time.Time) int {

	var econetDate int

	return econetDate
}

//int getEconetDate(dir_t* dirEntry) {
//int ecoDate;
//ecoDate = (FAT_YEAR(dirEntry->lastWriteDate) - 1981) << 1 & 224;
//ecoDate += FAT_DAY(dirEntry->lastWriteDate);
//ecoDate += (FAT_YEAR(dirEntry->lastWriteDate) - 1981) << 12;
//ecoDate += FAT_MONTH(dirEntry->lastWriteDate) << 8;
//
//}
