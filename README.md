# README

## Usage

Converting zonefile to Terraform
```
$ GLESYS_USERID=CL12345 GLESYS_TOKEN=ABC123FFF go run main.go -in=myzonefile -tf=dns-records.tf
```

Export ZoneFile from Glesys
```
$ GLESYS_USERID=CL12345 GLESYS_TOKEN=ABC123FFF go run main.go -domain=example.com -out=myzonefile
```

## Caution

Use at own risk. Make backups.
