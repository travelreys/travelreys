/*
   Every object has a real number as an index and the order of the children for an element
   of the tree is determined by sorting all children by their index. To insert between two objects,
   just set the index for the new object to the average index of the two objects on either side.
   We use arbitrary-precision fractions instead of 64-bit doubles so that we can’t run out of precision
   after lots of edits.

   Every index is a fraction between 0 and 1 exclusive. Being exclusive is important; it ensures we can always
   generate an index before or after an existing index by averaging with 0 or 1, respectively.
   Each index is stored as a string and averaging is done using string manipulation to retain precision.
   For compactness, we omit the leading “0.” from the fraction and we use the entire ASCII range instead of
   just the numbers 0–9 (base 95 instead of base 10).

   See:
    * https://www.figma.com/blog/realtime-editing-of-ordered-sequences
    * https://www.figma.com/blog/how-figmas-multiplayer-technology-works
    * https://observablehq.com/@dgreensp/implementing-fractional-indexing
*/


/*
  Fraction strings
  ----------------
  Represent the digits of a decimal as a string, so instead of 0.5, it's "5".
  When we switch from numbers to strings, the natural order we are taking advantage of is
  lexicographical order rather than numeric order.  For example, the string "2.1" sorts after
  the string "10.5".
    * Zero is the empty string, which sorts first.
    * To represent one, we can use a sentinel value such as END or null.
    * Ban trailing zeros, and treat the end of the string as equivalent to an infinite sequence
      of zero digits, exactly as we do when writing numbers.

*/

/*
  Midpoint between fractions
  --------------------------
  Now what we need is a "midpoint" algorithm that takes two string order keys and produces one in the middle.
  It doesn't have to be an exact midpoint, just roughly in the middle.
    * A is an order key string, or an empty string (START)
    * B is an order key string, greater than A, or null (END)
    * ...where an order key is a non-empty string of digits with no trailing zeros
    * Returns a key C such that A < C < B

  One property we'd like our midpoint algorithm to have is it doesn't needlessly grow the number of digits.
  We can exhaust the possibilities at the current digit length before adding a digit.

  Here is the algorithm for midpoint(A, B):

  1. If there is a non-empty common prefix of A and B, set it aside and recurse on the parts of A and B where they
     differ. For example, given "123" and "1234", call midpoint("", "4"), and then put the common prefix "123" back
     at the beginning, yielding "1232". When finding the longest common prefix, we must treat A as if it is followed
     by as many trailing zeros as we need to match zeros in B. For example, "123" and "123004" have a common prefix
     of "12300". If B is null, this step does not apply.
  2. If we get to this step, A and B must have different first digits. For the purposes of the following steps, if A
     is empty, we'll say the first digit is 0. If B is null, we'll say its first digit is 10 (an imaginary digit one
     larger than 9). We know the first digit of B is strictly greater than the first digit of A.
  3. If there are digits that are strictly between the first digit of A and the first digit of B, pick the median
     one and return that. For example, if A starts with "2" and B starts with "4", return "3".
  4. If we get to this step, B is one more than A. If the first digit of B would be a suitable return value, return that.
     For example, midpoint("35", "41") can return "4". If B is non-null and not just a single digit, its first digit is
     a suitable return value.
  5. If we get to this step, B is the smallest string that doesn't share a first digit with A, or null.
     For example, if A is "823", B is "9". If A is "923", B is null. We take the first digit off of A;
     take the midpoint of the rest with null; and put the first digit back on the beginning. So in the first
     example, "8" + midpoint("23", null). In the second example, "9" + midpoint("23", null).
*/

// Digits must be in ascending character code order!
const BASE_95_DIGITS = ' !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~';
const INTEGER_ZERO = "a0";
const SMALLEST_INTEGER = "A00000000000000000000000000";

// `a` may be empty string, `b` is null or non-empty string.
// `a < b` lexicographically if `b` is non-null.
// no trailing zeros allowed.
export const midpoint = (a: string, b: string | null): string => {
  const digits = BASE_95_DIGITS;

  if (b !== null && a >= b) {
    throw new Error(a + ' >= ' + b)
  }
  if (a.slice(-1) === '0' || (b && b.slice(-1) === '0')) {
    throw new Error('trailing zero')
  }

  if (b) {
    // remove longest common prefix.  pad `a` with 0s as we
    // go.  note that we don't need to pad `b`, because it can't
    // end before `a` while traversing the common prefix.
    let n = 0
    while ((a.charAt(n) || '0') === b.charAt(n)) {
      n++
    }
    if (n > 0) {
      return b.slice(0, n) + midpoint(a.slice(n), b.slice(n))
    }
  }

  // first digits (or lack of digit) are different
  const digitA = a ? digits.indexOf(a.charAt(0)) : 0
  const digitB = b !== null ? digits.indexOf(b.charAt(0)) : digits.length
  if (digitB - digitA > 1) {
    const midDigit = Math.round(0.5*(digitA + digitB))
    return digits.charAt(midDigit)
  } else {
    // first digits are consecutive
    if (b && b.length > 1) {
      return b.slice(0, 1)
    } else {
      // `b` is null or has length 1 (a single digit).
      // the first digit of `a` is the previous digit to `b`,
      // or 9 if `b` is null.
      // given, for example, midpoint('49', '5'), return
      // '4' + midpoint('9', null), which will become
      // '4' + '9' + midpoint('', null), which is '495'
      return digits.charAt(digitA) + midpoint(a.slice(1), null)
    }
  }
}

/*
  Logarithmic key growth
  ----------------------
  Now let's address the fact that even just adding items to the end of the list grows the length of the key at a
  linear rate. It would be nice to design a scheme where you can add to the end of the list with logarithmic key
  growth, and while we're at it, why not allow prepending with logarithmic key growth as well?

  Going back to numeric keys for a moment, what if we gave keys an integer part in addition to the fractional part?
  We could give list items keys like 1, 2, and 3, and only if you insert or move an item between these keys, create
  keys like 0.5, 1.5, and 2.5. When finding a "midpoint" between a key and END, we would always choose the smallest
  integer that is larger than the key.

  To represent 2.5 as a lexicographically-ordered string key, we can encode 2 as a string somehow, then encode 0.5 as
  "5", and then concatenate the strings. The immediate problem is how to get integers to sort in the correct order,
  since "2" is greater than "10" in lexicographical order.

  Encoding integers
  ------------------
  One solution is to pad the integer part to a fixed length. For example, we could decide on ten digits and count
  "0000000001", "0000000002", and so on as our order keys when appending items to an empty list. To illustrate how
  the fraction part comes in, the midpoint of "0000000001" and "0000000002" would be "00000000015". To support both
  prepending and appending to the list without going into the fractional part, we could start in the middle of the
  integers, with the first list item getting key "5000000000". Now appending and prepending actually don't change the
  length of the key at all, at least for a while. Once we run out of integers, we have to start using the fractional
  part to append or prepend new items. For example, after 9999999999 would come 9999999995. However, with base 95,
  ten digits gives us a little more than 65 bits of integer precision, which is many billions of billions of items.

  A variable-length encoding of integers that sorts correctly is also possible. We can prefix integers by their length.
  For example, 1 to 9 could be written A1 to A9, while 10 to 99 are written B10 to B99. In fact, we might as well use
  A0 to A9 and B00 to B99. In this scheme, the first item in an empty list gets key "a0". An item inserted after it
  gets "a1". An item inserted before it gets "Z9". The midpoint of "Z9" and "a0" is "Z95".

  Note that there is a smallest and largest integer in this scheme, if we go down to A00000000000000000000000000
  (A followed by 26 As) and up to zzzzzzzzzzzzzzzzzzzzzzzzzzz (z followed by 26 more zs). The set of integers is
  astronomical, but finite.

  Whichever integer system is used, the following operations are needed, and define an integer system:
    * Validate an integer
    * Extract an integer from the beginning of a string
    * Get zero
    * Increment an integer, returning null if it is the largest available integer
    * Decrement an integer, returning null if it is the smallest available integer
    * As a convenience, get the smallest integer (though an integer can be tested to see if it is
      the smallest by trying to decrement it)
*/

const getIntegerLength = (head: string): number => {
  if (head >= 'a' && head <= 'z') {
    return head.charCodeAt(0) - 'a'.charCodeAt(0) + 2
  } else if (head >= 'A' && head <= 'Z') {
    return 'Z'.charCodeAt(0) - head.charCodeAt(0) + 2
  } else {
    throw new Error('Invalid order key head: ' + head)
  }
}

const validateInteger = (int: string) => {
  if (int.length !== getIntegerLength(int.charAt(0))) {
     throw new Error('invalid integer part of order key: ' + int)
  }
}

const incrementInteger = (x: string): string|null => {
  validateInteger(x)
  const [head, ...digs] = x.split('')
  let carry = true
  for (let i = digs.length - 1; carry && i >= 0; i--) {
    const d = BASE_95_DIGITS.indexOf(digs[i]) + 1
    if (d === BASE_95_DIGITS.length) {
      digs[i] = '0'
    } else {
      digs[i] = BASE_95_DIGITS.charAt(d)
      carry = false
    }
  }
  if (carry) {
    if (head === 'Z') {
      return 'a0'
    }
    if (head === 'z') {
      return null
    }
    const h = String.fromCharCode(head.charCodeAt(0) + 1)
    if (h > 'a') {
      digs.push('0')
    } else {
      digs.pop()
    }
    return h + digs.join('')
  } else {
    return head + digs.join('')
  }
}

const decrementInteger = (x: string): string|null => {
  validateInteger(x)
  const [head, ...digs] = x.split('')
  let borrow = true
  for (let i = digs.length - 1; borrow && i >= 0; i--) {
    const d = BASE_95_DIGITS.indexOf(digs[i]) - 1
    if (d === -1) {
      digs[i] = BASE_95_DIGITS.slice(-1)
    } else {
      digs[i] = BASE_95_DIGITS.charAt(d)
      borrow = false
    }
  }
  if (borrow) {
    if (head === 'a') {
      return 'Z' + BASE_95_DIGITS.slice(-1)
    }
    if (head === 'A') {
      return null
    }
    const h = String.fromCharCode(head.charCodeAt(0) - 1)
    if (h < 'Z') {
      digs.push(BASE_95_DIGITS.slice(-1))
    } else {
      digs.pop()
    }
    return h + digs.join('')
  } else {
    return head + digs.join('')
  }
}

/*
  Algorithm combining integer and fraction encoding
  -------------------------------------------------
  An order key is now a string which can easily and unambiguously be divided into two strings, an integer part followed
  by a fraction part. The integer part will never be an empty string, and the fraction part will never be null. Instead
  of using the empty string as START and null as END, we'll refer explicitly to START and END, and use null to mean
  START or END when we need a sentinel value in JavaScript. Note that the empty string was not a valid order key before,
  but it is a valid fraction now, but only as long as the integer part is not equal to the smallest integer. We can't
  allow there to be an order key such that we can't construct a smaller order key, and this would happen if we allowed
  the smallest available integer as a key with no fraction part. (Even though this case will not be reached by prepending
  items to a list in normal operation because it would require an astronomical number of items, we don't want any
  undefined edge cases.)

  Ok, so in terms of preconditions and postconditions, we have:
  generateKeyBetween(A, B)
    * A is an order key string, or START
    * B is an order key string, greater than A, or END
    * An order key string is an integer part followed by a fraction part
    * The fraction part consists of zero or more digits with no trailing zeros
    * If the fraction part is empty, the integer part may not be the smallest available integer
    * Returns a key C such that A < C < B

  The algorithm will use the fraction midpoint algorithm unchanged, as well as the capability to increment
  and decrement integers. There is no need to actually convert our "string integers" into numeric integers or
  vice versa. We just need a "zero" and a way to increment and decrement.

  Algorithm:
  * If A is START and B is END, return the "zero" integer with an empty fraction part.
  * If A is START and the integer part of B is the smallest available integer, make the fraction part of B smaller
    using the midpoint algorithm with the empty string as the first argument. Otherwise, if A is START, return the
    largest integer less than B, with an empty fraction part.
  * If B is END, return the smallest integer greater than A, with an empty fraction part. If there is no such integer,
    make the fraction part of A larger, using the midpoint algorithm with a second argument of null.
  * If A and B have the same integer part, find the midpoint of the fraction parts, and return a key consisting of the
    common integer part and the midpoint.
  * Increment the integer part of A. (This is always possible, because B has a larger integer part.) If a key
    consisting of this integer and no fraction is less than B, return it. In other words, if A is the string
    equivalent of the number 3.5, see if we can return 4.
  * Otherwise, take the integer part of A and make the fraction part larger, using the midpoint algorithm with a
    second argument of null.
*/

const getIntegerPart = (key: string): string =>  {
  const integerPartLength = getIntegerLength(key.charAt(0))
  if (integerPartLength > key.length) {
    throw new Error('invalid order key: ' + key)
  }
  return key.slice(0, integerPartLength)
}

const validateOrderKey = (key: string) => {
  if (key === SMALLEST_INTEGER) {
    throw new Error('invalid order key: ' + key)
  }
  // getIntegerPart will throw if the first character is bad,
  // or the key is too short.  we'd call it to check these things
  // even if we didn't need the result
  const i = getIntegerPart(key)
  const f = key.slice(i.length)
  if (f.slice(-1) === '0') {
    throw new Error('invalid order key: ' + key)
  }
}

// `a` is an order key or null (START).
// `b` is an order key or null (END).
// `a < b` lexicographically if both are non-null.
export const generateKeyBetween = (a: string|null, b: string|null): any => {
  if (a !== null) {
    validateOrderKey(a)
  }
  if (b !== null) {
    validateOrderKey(b)
  }
  if (a !== null && b !== null && a >= b) {
    throw new Error(a + ' >= ' + b)
  }
  if (a === null && b === null) {
    return INTEGER_ZERO
  }

  b = b as string

  if (a === null) {
    const ib = getIntegerPart(b as string)
    const fb = (b as string).slice(ib.length)
    if (ib === SMALLEST_INTEGER) {
      return ib + midpoint('', fb)
    }
    return ib < b ? ib : decrementInteger(ib)
  }
  if (b === null) {
    const ia = getIntegerPart(a)
    const fa = a.slice(ia.length)
    const i = incrementInteger(ia)
    return i === null ? ia + midpoint(fa, null) : i
  }
  const ia = getIntegerPart(a)
  const fa = a.slice(ia.length)
  const ib = getIntegerPart(b)
  const fb = b.slice(ib.length)
  if (ia === ib) {
    return ia + midpoint(fa, fb)
  }
  const i = incrementInteger(ia)
  return i! < b ? i : ia + midpoint(fa, null)
}
