CREATE TABLE IF NOT EXISTS good_logs
(
    Id          Int64,
    ProjectId   Int64,
    Name        String,
    Description String,
    Priority    Int,
    Removed     Boolean,
    EventTime   DateTime
) ENGINE = MergeTree()
      ORDER BY (EventTime);