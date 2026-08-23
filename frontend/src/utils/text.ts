function pluralize(
  count: number,
  singular: string,
  plural: string,
  includeCount: boolean,
): string {
  const word = count === 1 ? singular : plural;
  return includeCount ? `${count} ${word}` : word;
}

export { pluralize };
