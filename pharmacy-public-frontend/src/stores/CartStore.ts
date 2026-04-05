import { defineStore } from "pinia";
import type { CartState } from "@/types";

export const useCartStore = defineStore("cart", {
  state: (): CartState => ({
    items: JSON.parse(localStorage.getItem("cart") || "{}")
  }),
  actions: {
    addToCart(productId: number) {
      if (this.items[productId]) {
        this.items[productId]++;
      } else {
        this.items[productId] = 1;
      }
      localStorage.setItem("cart", JSON.stringify(this.items));
    },
    removeFromCart(productId: number) {
      if (this.items[productId]) {
        if (this.items[productId] > 1) {
          this.items[productId]--;
        } else {
          delete this.items[productId];
        }
      }
      localStorage.setItem("cart", JSON.stringify(this.items));
    },
    clearCart() {
      this.items = {};
      localStorage.setItem("cart", JSON.stringify(this.items));
    }
  }
});



