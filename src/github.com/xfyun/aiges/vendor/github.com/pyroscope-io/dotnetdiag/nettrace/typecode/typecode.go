package typecode

type TypeCode int32

// https://docs.microsoft.com/en-us/dotnet/api/system.typecode?view=net-5.0
const (
	_ TypeCode = iota - 1

	Empty    // A null reference.
	Object   // A general type representing any reference or value type not explicitly represented by another TypeCode.
	DBNull   // A database null (column) value.
	Boolean  // A simple type representing Boolean values of true or false.
	Char     // An integral type representing unsigned 16-bit integers with values between 0 and 65535. The set of possible values for the Char type corresponds to the Unicode character set.
	SByte    // An integral type representing signed 8-bit integers with values between -128 and 127.
	Byte     // An integral type representing unsigned 8-bit integers with values between 0 and 255.
	Int16    // An integral type representing signed 16-bit integers with values between -32768 and 32767.
	UInt16   // An integral type representing unsigned 16-bit integers with values between 0 and 65535.
	Int32    // An integral type representing signed 32-bit integers with values between -2147483648 and 2147483647.
	UInt32   // An integral type representing unsigned 32-bit integers with values between 0 and 4294967295.
	Int64    // An integral type representing signed 64-bit integers with values between -9223372036854775808 and 9223372036854775807.
	UInt64   // An integral type representing unsigned 64-bit integers with values between 0 and 18446744073709551615.
	Single   // A floating point type representing values ranging from approximately 1.5 x 10 -45 to 3.4 x 10 38 with a precision of 7 digits.
	Double   // A floating point type representing values ranging from approximately 5.0 x 10 -324 to 1.7 x 10 308 with a precision of 15-16 digits.
	Decimal  // A simple type representing values ranging from 1.0 x 10 -28 to approximately 7.9 x 10 28 with 28-29 significant digits.
	DateTime // A type representing a date and time value.
	_        // http://blogs.ugidotnet.org/adrian/archive/2008/03/17/91755.aspx
	String   // A sealed class type representing Unicode character strings.
	Array    // https://github.com/microsoft/perfview/blob/main/src/TraceEvent/EventPipe/EventPipeFormat.md#metadata-event-encoding
)
