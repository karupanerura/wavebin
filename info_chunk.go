package wavebin

import (
	"github.com/karupanerura/riffbin"
)

type InfoKey [4]byte

var (
	InfoRatedAGES               = InfoKey{'A', 'G', 'E', 'S'}
	InfoCommentCMNT             = InfoKey{'C', 'M', 'N', 'T'}
	InfoEncodedByCODE           = InfoKey{'C', 'O', 'D', 'E'}
	InfoCommentsCOMM            = InfoKey{'C', 'O', 'M', 'M'}
	InfoDirectoryDIRC           = InfoKey{'D', 'I', 'R', 'C'}
	InfoSoundSchemeTitleDISP    = InfoKey{'D', 'I', 'S', 'P'}
	InfoDateTimeOriginalDTIM    = InfoKey{'D', 'T', 'I', 'M'}
	InfoGenreGENR               = InfoKey{'G', 'E', 'N', 'R'}
	InfoArchivalLocationIARL    = InfoKey{'I', 'A', 'R', 'L'}
	InfoArtistIART              = InfoKey{'I', 'A', 'R', 'T'}
	InfoFirstLanguageIAS1       = InfoKey{'I', 'A', 'S', '1'}
	InfoSecondLanguageIAS2      = InfoKey{'I', 'A', 'S', '2'}
	InfoThirdLanguageIAS3       = InfoKey{'I', 'A', 'S', '3'}
	InfoFourthLanguageIAS4      = InfoKey{'I', 'A', 'S', '4'}
	InfoFifthLanguageIAS5       = InfoKey{'I', 'A', 'S', '5'}
	InfoSixthLanguageIAS6       = InfoKey{'I', 'A', 'S', '6'}
	InfoSeventhLanguageIAS7     = InfoKey{'I', 'A', 'S', '7'}
	InfoEighthLanguageIAS8      = InfoKey{'I', 'A', 'S', '8'}
	InfoNinthLanguageIAS9       = InfoKey{'I', 'A', 'S', '9'}
	InfoBaseURLIBSU             = InfoKey{'I', 'B', 'S', 'U'}
	InfoDefaultAudioStreamICAS  = InfoKey{'I', 'C', 'A', 'S'}
	InfoCostumeDesignerICDS     = InfoKey{'I', 'C', 'D', 'S'}
	InfoCommissionedICMS        = InfoKey{'I', 'C', 'M', 'S'}
	InfoCommentICMT             = InfoKey{'I', 'C', 'M', 'T'}
	InfoCinematographerICNM     = InfoKey{'I', 'C', 'N', 'M'}
	InfoCountryICNT             = InfoKey{'I', 'C', 'N', 'T'}
	InfoCopyrightICOP           = InfoKey{'I', 'C', 'O', 'P'}
	InfoDateCreatedICRD         = InfoKey{'I', 'C', 'R', 'D'}
	InfoCroppedICRP             = InfoKey{'I', 'C', 'R', 'P'}
	InfoDimensionsIDIM          = InfoKey{'I', 'D', 'I', 'M'}
	InfoDateTimeOriginalIDIT    = InfoKey{'I', 'D', 'I', 'T'}
	InfoDotsPerInchIDPI         = InfoKey{'I', 'D', 'P', 'I'}
	InfoDistributedByIDST       = InfoKey{'I', 'D', 'S', 'T'}
	InfoEditedByIEDT            = InfoKey{'I', 'E', 'D', 'T'}
	InfoEncodedByIENC           = InfoKey{'I', 'E', 'N', 'C'}
	InfoEngineerIENG            = InfoKey{'I', 'E', 'N', 'G'}
	InfoGenreIGNR               = InfoKey{'I', 'G', 'N', 'R'}
	InfoKeywordsIKEY            = InfoKey{'I', 'K', 'E', 'Y'}
	InfoLightnessILGT           = InfoKey{'I', 'L', 'G', 'T'}
	InfoLogoURLILGU             = InfoKey{'I', 'L', 'G', 'U'}
	InfoLogoIconURLILIU         = InfoKey{'I', 'L', 'I', 'U'}
	InfoLanguageILNG            = InfoKey{'I', 'L', 'N', 'G'}
	InfoMoreInfoBannerImageIMBI = InfoKey{'I', 'M', 'B', 'I'}
	InfoMoreInfoBannerURLIMBU   = InfoKey{'I', 'M', 'B', 'U'}
	InfoMediumIMED              = InfoKey{'I', 'M', 'E', 'D'}
	InfoMoreInfoTextIMIT        = InfoKey{'I', 'M', 'I', 'T'}
	InfoMoreInfoURLIMIU         = InfoKey{'I', 'M', 'I', 'U'}
	InfoMusicByIMUS             = InfoKey{'I', 'M', 'U', 'S'}
	InfoTitleINAM               = InfoKey{'I', 'N', 'A', 'M'}
	InfoProductionDesignerIPDS  = InfoKey{'I', 'P', 'D', 'S'}
	InfoNumColorsIPLT           = InfoKey{'I', 'P', 'L', 'T'}
	InfoProductIPRD             = InfoKey{'I', 'P', 'R', 'D'}
	InfoProducedByIPRO          = InfoKey{'I', 'P', 'R', 'O'}
	InfoRippedByIRIP            = InfoKey{'I', 'R', 'I', 'P'}
	InfoRatingIRTD              = InfoKey{'I', 'R', 'T', 'D'}
	InfoSubjectISBJ             = InfoKey{'I', 'S', 'B', 'J'}
	InfoSoftwareISFT            = InfoKey{'I', 'S', 'F', 'T'}
	InfoSecondaryGenreISGN      = InfoKey{'I', 'S', 'G', 'N'}
	InfoSharpnessISHP           = InfoKey{'I', 'S', 'H', 'P'}
	InfoTimeCodeISMP            = InfoKey{'I', 'S', 'M', 'P'}
	InfoSourceISRC              = InfoKey{'I', 'S', 'R', 'C'}
	InfoSourceFormISRF          = InfoKey{'I', 'S', 'R', 'F'}
	InfoProductionStudioISTD    = InfoKey{'I', 'S', 'T', 'D'}
	InfoStarringISTR            = InfoKey{'I', 'S', 'T', 'R'}
	InfoTechnicianITCH          = InfoKey{'I', 'T', 'C', 'H'}
	InfoTrackNumberITRK         = InfoKey{'I', 'T', 'R', 'K'}
	InfoWatermarkURLIWMU        = InfoKey{'I', 'W', 'M', 'U'}
	InfoWrittenByIWRI           = InfoKey{'I', 'W', 'R', 'I'}
	InfoLanguageLANG            = InfoKey{'L', 'A', 'N', 'G'}
	InfoLocationLOCA            = InfoKey{'L', 'O', 'C', 'A'}
	InfoPartPRT1                = InfoKey{'P', 'R', 'T', '1'}
	InfoNumberOfPartsPRT2       = InfoKey{'P', 'R', 'T', '2'}
	InfoRateRATE                = InfoKey{'R', 'A', 'T', 'E'}
	InfoStarringSTAR            = InfoKey{'S', 'T', 'A', 'R'}
	InfoStatisticsSTAT          = InfoKey{'S', 'T', 'A', 'T'}
	InfoTapeNameTAPE            = InfoKey{'T', 'A', 'P', 'E'}
	InfoEndTimecodeTCDO         = InfoKey{'T', 'C', 'D', 'O'}
	InfoStartTimecodeTCOD       = InfoKey{'T', 'C', 'O', 'D'}
	InfoTitleTITL               = InfoKey{'T', 'I', 'T', 'L'}
	InfoLengthTLEN              = InfoKey{'T', 'L', 'E', 'N'}
	InfoOrganizationTORG        = InfoKey{'T', 'O', 'R', 'G'}
	InfoTrackNumberTRCK         = InfoKey{'T', 'R', 'C', 'K'}
	InfoURLTURL                 = InfoKey{'T', 'U', 'R', 'L'}
	InfoVersionTVER             = InfoKey{'T', 'V', 'E', 'R'}
	InfoVegasVersionMajorVMAJ   = InfoKey{'V', 'M', 'A', 'J'}
	InfoVegasVersionMinorVMIN   = InfoKey{'V', 'M', 'I', 'N'}
	InfoYearYEAR                = InfoKey{'Y', 'E', 'A', 'R'}
)

type InfoChunk struct {
	Data map[InfoKey]string
}

func (f *InfoChunk) Chunk() riffbin.Chunk {
	payload := make([]riffbin.Chunk, 0, len(f.Data))
	for key, value := range f.Data {
		payload = append(payload, &riffbin.OnMemorySubChunk{
			ID:      key,
			Payload: []byte(value),
		})
	}

	return &riffbin.ListChunk{
		ListType: infoBytes,
		Payload:  payload,
	}
}
