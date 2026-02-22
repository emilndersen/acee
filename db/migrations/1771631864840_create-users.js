exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable("users", {
    id: "id",
    name: { type: "text", notNull: true },
    email: { type: "text", notNull: true, unique: true },
    created_at: {
      type: "timestamptz",
      notNull: true,
      default: pgm.func("now()"),
    },
  });
};

exports.down = (pgm) => {
  pgm.dropTable("users");
};