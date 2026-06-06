interface Credentials {
  username: string;
  password: string;
  displayName: string;
}

type PlayerInfoData = {
  displayName: string;
  username: string;
  isAuthenticated: boolean;
};

export type { Credentials, PlayerInfoData };
