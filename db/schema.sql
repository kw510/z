CREATE TABLE IF NOT EXISTS tuples (
  parent_namespace TEXT NOT NULL,
  parent_id TEXT NOT NULL,
  parent_relation TEXT NOT NULL,
  child_namespace TEXT NOT NULL,
  child_id TEXT NOT NULL,
  child_relation TEXT NOT NULL,

  PRIMARY KEY (parent_namespace, parent_id, parent_relation, child_namespace, child_id, child_relation)
);

CREATE INDEX IF NOT EXISTS idx_tuples_parent ON tuples (parent_namespace, parent_id, parent_relation);
CREATE INDEX IF NOT EXISTS idx_tuples_child ON tuples (child_namespace, child_id, child_relation);

CREATE TABLE IF NOT EXISTS relations (
  namespace TEXT NOT NULL,
  relation TEXT NOT NULL,
  permission TEXT NOT NULL,

  PRIMARY KEY (namespace, relation, permission)
);

CREATE INDEX IF NOT EXISTS idx_relations_namespace ON relations (namespace);
CREATE INDEX IF NOT EXISTS idx_relations_relation ON relations (relation);
CREATE INDEX IF NOT EXISTS idx_relations_permission ON relations (permission);
