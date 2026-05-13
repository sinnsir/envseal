// Package store provides a simple file-based persistence layer for sealed
// environment files produced by envseal.
//
// Sealed files are stored under a configurable base directory (default:
// .envseal/) with the naming convention <environment>.sealed, e.g.:
//
//	.envseal/
//	  production.sealed
//	  staging.sealed
//	  development.sealed
//
// Each file contains the raw bytes returned by [envelope.Marshal] and is
// safe to commit to version control — the contents are age-encrypted and
// carry no plaintext secrets.
package store
