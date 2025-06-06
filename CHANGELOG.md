## 0.1.2
* Add Float64 type to handle Location attributes; they are sometimes returned as floats and sometimes as strings

## 0.1.1
 * Fix http err handling; add tests

## 0.1.0

Initial release. Supports CRUD(-ish) ops for these resources (primarily for [terraform-provider-librenms](https://github.com/jokelyo/terraform-provider-librenms)):
 * Alert Rules
 * Devices
 * Device Groups
 * Locations
 * Services
