import { Filter } from "./filter";
import { Sort } from "./sort";

export interface PageRequest {
  page: number;

  size: number;

  filters?: Filter[];

  sorts?: Sort[];
}
