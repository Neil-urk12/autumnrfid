/**
 * Formats a numeric year level into a human-readable string.
 * e.g., 1 -> "1st Year", 2 -> "2nd Year", null -> "N/A"
 * @param yearLevel The numeric year level or null.
 * @returns The formatted string representation.
 */
export function formatYearLevel(yearLevel: number | null): string {
  if (yearLevel === null || yearLevel === undefined || yearLevel <= 0) {
    return 'N/A'; // Or handle as appropriate (e.g., empty string, 'Unknown')
  }

  switch (yearLevel) {
    case 1:
      return '1st Year';
    case 2:
      return '2nd Year';
    case 3:
      return '3rd Year';
    default:
      return `${yearLevel}th Year`;
  }
}
