import { formatMatchTime, formatTimeFormat } from "../utils/time";

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
