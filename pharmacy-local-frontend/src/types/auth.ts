export interface User {
  id: number;
  username: string;
  permissions?: String[];
}

export interface AuthResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

export interface LoginData {
  login: string;
  password: string;
}
