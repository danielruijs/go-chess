import {
  formatMatchTime,
  formatTimeFormat,
  formatJoinDate,
  formatMatchDate,
  formatBaseTime,
  formatIncrement,
} from "../utils/time";

describe("formatMatchTime", () => {
  it("formats 0ms as 0:00", () => {
    expect(formatMatchTime(0)).toBe("0:00");
  });

  it("formats 1ms as 0:00.0", () => {
    expect(formatMatchTime(1)).toBe("0:00.0");
  });

  it("formats 18ms as 0:00.0", () => {
    expect(formatMatchTime(18)).toBe("0:00.0");
  });

  it("formats 350ms as 0:00.3", () => {
    expect(formatMatchTime(350)).toBe("0:00.3");
  });

  it("formats 497ms as 0:00.4", () => {
    expect(formatMatchTime(497)).toBe("0:00.4");
  });

  it("formats 500ms as 0:00.5", () => {
    expect(formatMatchTime(500)).toBe("0:00.5");
  });

  it("formats 900ms as 0:00.9", () => {
    expect(formatMatchTime(900)).toBe("0:00.9");
  });

  it("formats 1050ms as 0:01.0", () => {
    expect(formatMatchTime(1050)).toBe("0:01.0");
  });

  it("formats 1500ms as 0:01.5", () => {
    expect(formatMatchTime(1500)).toBe("0:01.5");
  });

  it("formats 3782ms as 0:03.7", () => {
    expect(formatMatchTime(3782)).toBe("0:03.7");
  });

  it("formats 9000ms as 0:09.0", () => {
    expect(formatMatchTime(9000)).toBe("0:09.0");
  });

  it("formats 11 500ms as 0:11", () => {
    expect(formatMatchTime(11500)).toBe("0:11");
  });

  it("formats 59 000ms as 0:59", () => {
    expect(formatMatchTime(59 * 1000)).toBe("0:59");
  });

  it("formats 59999ms as 0:59", () => {
    expect(formatMatchTime(59999)).toBe("0:59");
  });

  it("formats 60 000ms as 1:00", () => {
    expect(formatMatchTime(60 * 1000)).toBe("1:00");
  });

  it("formats 125 000ms as 2:05", () => {
    expect(formatMatchTime(125 * 1000)).toBe("2:05");
  });

  it("formats 300 000ms as 5:00", () => {
    expect(formatMatchTime(300 * 1000)).toBe("5:00");
  });
});

describe("formatTimeFormat", () => {
  it("formats 3min initial and 0s increment as 3+0", () => {
    expect(formatTimeFormat({ initialMs: 180 * 1000, incrementMs: 0 })).toBe("3+0");
  });

  it("formats 1min initial and 1s increment as 1+1", () => {
    expect(formatTimeFormat({ initialMs: 60 * 1000, incrementMs: 1000 })).toBe("1+1");
  });

  it("formats 15min initial and 10s increment as 15+10", () => {
    expect(formatTimeFormat({ initialMs: 15 * 60 * 1000, incrementMs: 10 * 1000 })).toBe("15+10");
  });
});

describe("formatBaseTime", () => {
  it("formats 1min as 1 minute", () => {
    expect(formatBaseTime({ initialMs: 60 * 1000, incrementMs: 0 })).toBe("1 minute");
  });

  it("formats 3min as 3 minutes", () => {
    expect(formatBaseTime({ initialMs: 180 * 1000, incrementMs: 2000 })).toBe("3 minutes");
  });

  it("formats 15min as 15 minutes", () => {
    expect(formatBaseTime({ initialMs: 15 * 60 * 1000, incrementMs: 10 * 1000 })).toBe(
      "15 minutes"
    );
  });

  it("formats 30s as 30 seconds", () => {
    expect(formatBaseTime({ initialMs: 30 * 1000, incrementMs: 0 })).toBe("30 seconds");
  });

  it("formats 1s as 1 second", () => {
    expect(formatBaseTime({ initialMs: 1000, incrementMs: 0 })).toBe("1 second");
  });

  it("formats 90s as 1m 30s", () => {
    expect(formatBaseTime({ initialMs: 90 * 1000, incrementMs: 0 })).toBe("1m 30s");
  });

  it("formats 0ms as 0 seconds", () => {
    expect(formatBaseTime({ initialMs: 0, incrementMs: 0 })).toBe("0 seconds");
  });
});

describe("formatIncrement", () => {
  it("formats 0s increment as None", () => {
    expect(formatIncrement({ initialMs: 60 * 1000, incrementMs: 0 })).toBe("None");
  });

  it("formats 1s increment as 1 second", () => {
    expect(formatIncrement({ initialMs: 60 * 1000, incrementMs: 1000 })).toBe("1 second");
  });

  it("formats 2s increment as 2 seconds", () => {
    expect(formatIncrement({ initialMs: 180 * 1000, incrementMs: 2000 })).toBe("2 seconds");
  });

  it("formats 10s increment as 10 seconds", () => {
    expect(formatIncrement({ initialMs: 15 * 60 * 1000, incrementMs: 10 * 1000 })).toBe(
      "10 seconds"
    );
  });
});

describe("formatJoinDate", () => {
  it("formats valid date string", () => {
    const result = formatJoinDate("2024-01-15T12:00:00Z");
    expect(result).toContain("2024");
    expect(result).toContain("January");
  });

  it("falls back to raw string on invalid date", () => {
    expect(formatJoinDate("invalid-date")).toBe("invalid-date");
  });
});

describe("formatMatchDate", () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.setSystemTime(new Date("2026-08-16T12:00:00Z"));
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it("formats less than 1 min as Just now", () => {
    const recent = "2026-08-16T11:59:45Z";
    expect(formatMatchDate(recent)).toBe("Just now");
  });

  it("formats minutes ago", () => {
    const minAgo = "2026-08-16T11:45:00Z";
    expect(formatMatchDate(minAgo)).toBe("15m ago");
  });

  it("formats hours ago", () => {
    const hoursAgo = "2026-08-16T08:00:00Z";
    expect(formatMatchDate(hoursAgo)).toBe("4h ago");
  });

  it("formats days ago", () => {
    const daysAgo = "2026-08-13T12:00:00Z";
    expect(formatMatchDate(daysAgo)).toBe("3d ago");
  });

  it("falls back to formatted calendar date after 7 days without year for same year", () => {
    const oldDate = "2026-08-01T12:00:00Z";
    const result = formatMatchDate(oldDate);
    expect(result).toContain("Aug");
    expect(result).not.toContain("2026");
  });

  it("includes the year for a date in a previous year", () => {
    const prevYearDate = "2025-11-20T12:00:00Z";
    const result = formatMatchDate(prevYearDate);
    expect(result).toContain("Nov");
    expect(result).toContain("2025");
  });
});
