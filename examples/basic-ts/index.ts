import * as pulumi from "@pulumi/pulumi";
import * as cockroachsql from "@pulumi/cockroachsql";

const resource = new cockroachsql.Resource("Resource", { sampleAttribute: "attr" });

export const sampleAttribute = resource.sampleAttribute;
