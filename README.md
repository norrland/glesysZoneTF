# README

## Usage

Converting zonefile to Terraform format. [TF Docs: glesys_dnsdomain_record](https://registry.terraform.io/providers/glesys/glesys/latest/docs/resources/dnsdomain_record)
```
$ go run main.go -in=myzonefile -tf=dns-records.tf
```

Export ZoneFile from Glesys DNS
```
$ GLESYS_USERID=CL12345 GLESYS_TOKEN=ABC123FFF go run main.go -domain=example.com -out=myzonefile
```

## Caution

Use at own risk. Make backups.
