import { defineStore } from "pinia";
import axios from "axios";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    errors: {
      username: "",
      password: "",
      twoFA: "",
      general: "",
    },
    authToken: localStorage.getItem("wira_token") || null, // Fetch token from local storage
  }),
  actions: {
    async login(credentials) {
      try {
        // Create URL-encoded form data
        const formData = new URLSearchParams();
        formData.append("username", credentials.username);
        formData.append("password", credentials.password);
        formData.append("twoFA_code", credentials.twoFA_code);

        const response = await axios.post(
          "http://localhost:8080/login",
          formData.toString(),
          {
            headers: {
              "Content-Type": "application/x-www-form-urlencoded",
            },
          }
        );

        // Save token and clear errors
        localStorage.setItem("wira_token", response.data.token);
        this.authToken = response.data.token;
        this.clearErrors();
        return true;
      } catch (error) {
        this.handleError(error);
        throw error;
      }
    },
    logout() {
      // Clear token and redirect to login
      this.authToken = null;
      localStorage.removeItem("wira_token");
    },
    handleError(error) {
      this.clearErrors();
      const message = error.response?.data?.error || "Login failed";

      if (
        !message.includes("username") &&
        !message.includes("password") &&
        !message.includes("2FA")
      ) {
        this.errors.general = message;
      } else {
        if (message.includes("username")) this.errors.username = message;
        if (message.includes("password")) this.errors.password = message;
        if (message.includes("2FA")) this.errors.twoFA = message;
      }
    },
    clearErrors() {
      this.errors = {
        username: "",
        password: "",
        twoFA: "",
        general: "",
      };
    },
  },
});
