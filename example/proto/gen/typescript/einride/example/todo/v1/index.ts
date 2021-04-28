// Example: 'tenants/{tenant}/users/{user}/todos/{todo}'
export interface TodoResourceName {
  tenant: string;
  user: string;
  todo: string;
  userResourceName(): UserResourceName;
  tenantResourceName(): TenantResourceName;
  toString(): string;
}

interface TodoResourceNameConstructor {
  parse(s: string): TodoResourceName;
  from(tenant: string, user: string, todo: string): TodoResourceName;
}

export const TodoResourceName: TodoResourceNameConstructor = {
  parse(s: string): TodoResourceName {
    const errPrefix = `parse resource name ${s} as todo-example.einride.tech/Todo`;
    const segments = s.split("/")
    if (segments.length !== 6) {
      throw new Error(`${errPrefix}: invalid segment count ${segments.length} (expected 6)`)
    }
    if (segments[0] !== "tenants") {
      throw new Error(`${errPrefix}: invalid constant segment ${segments[0]} (expected tenants)`)
    }
    const tenant = segments[1]
    if (segments[2] !== "users") {
      throw new Error(`${errPrefix}: invalid constant segment ${segments[2]} (expected users)`)
    }
    const user = segments[3]
    if (segments[4] !== "todos") {
      throw new Error(`${errPrefix}: invalid constant segment ${segments[4]} (expected todos)`)
    }
    const todo = segments[5]
    return this.from(tenant, user, todo)
  },

  from(tenant: string, user: string, todo: string): TodoResourceName {
    if (tenant === "" || tenant.indexOf("/") > -1) {
      throw new Error(`invalid variable segment for tenant: '${tenant}'`)
    }
    if (user === "" || user.indexOf("/") > -1) {
      throw new Error(`invalid variable segment for user: '${user}'`)
    }
    if (todo === "" || todo.indexOf("/") > -1) {
      throw new Error(`invalid variable segment for todo: '${todo}'`)
    }
    return {
      tenant,
      user,
      todo,
      userResourceName(): UserResourceName {
        return UserResourceName.from(tenant, user)
      },
      tenantResourceName(): TenantResourceName {
        return TenantResourceName.from(tenant)
      },
      toString(): string {
        // eslint-disable-next-line no-useless-concat
        return "tenants" + "/" + tenant + "/" + "users" + "/" + user + "/" + "todos" + "/" + todo
      },
    }
  },
}

// Example: 'tenants/{tenant}/users/{user}'
interface UserResourceName {
  tenant: string;
  user: string;
  tenantResourceName(): TenantResourceName;
  toString(): string;
}

interface UserResourceNameConstructor {
  parse(s: string): UserResourceName;
  from(tenant: string, user: string): UserResourceName;
}

const UserResourceName: UserResourceNameConstructor = {
  parse(s: string): UserResourceName {
    const errPrefix = `parse resource name ${s} as account-example.einride.tech/User`;
    const segments = s.split("/")
    if (segments.length !== 4) {
      throw new Error(`${errPrefix}: invalid segment count ${segments.length} (expected 4)`)
    }
    if (segments[0] !== "tenants") {
      throw new Error(`${errPrefix}: invalid constant segment ${segments[0]} (expected tenants)`)
    }
    const tenant = segments[1]
    if (segments[2] !== "users") {
      throw new Error(`${errPrefix}: invalid constant segment ${segments[2]} (expected users)`)
    }
    const user = segments[3]
    return this.from(tenant, user)
  },

  from(tenant: string, user: string): UserResourceName {
    if (tenant === "" || tenant.indexOf("/") > -1) {
      throw new Error(`invalid variable segment for tenant: '${tenant}'`)
    }
    if (user === "" || user.indexOf("/") > -1) {
      throw new Error(`invalid variable segment for user: '${user}'`)
    }
    return {
      tenant,
      user,
      tenantResourceName(): TenantResourceName {
        return TenantResourceName.from(tenant)
      },
      toString(): string {
        // eslint-disable-next-line no-useless-concat
        return "tenants" + "/" + tenant + "/" + "users" + "/" + user
      },
    }
  },
}

// Example: 'tenants/{tenant}'
interface TenantResourceName {
  tenant: string;
  toString(): string;
}

interface TenantResourceNameConstructor {
  parse(s: string): TenantResourceName;
  from(tenant: string): TenantResourceName;
}

const TenantResourceName: TenantResourceNameConstructor = {
  parse(s: string): TenantResourceName {
    const errPrefix = `parse resource name ${s} as account-example.einride.tech/Tenant`;
    const segments = s.split("/")
    if (segments.length !== 2) {
      throw new Error(`${errPrefix}: invalid segment count ${segments.length} (expected 2)`)
    }
    if (segments[0] !== "tenants") {
      throw new Error(`${errPrefix}: invalid constant segment ${segments[0]} (expected tenants)`)
    }
    const tenant = segments[1]
    return this.from(tenant)
  },

  from(tenant: string): TenantResourceName {
    if (tenant === "" || tenant.indexOf("/") > -1) {
      throw new Error(`invalid variable segment for tenant: '${tenant}'`)
    }
    return {
      tenant,
      toString(): string {
        // eslint-disable-next-line no-useless-concat
        return "tenants" + "/" + tenant
      },
    }
  },
}

